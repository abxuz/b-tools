package brpc

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	"sync"
	"time"

	"github.com/abxuz/b-tools/bcrypt"
	"github.com/vmihailenco/msgpack/v5"
)

type internalError error

type serverCodec struct {
	RequestReader  io.Reader
	ResponseWriter io.Writer
}

func (c *serverCodec) ReadRequestHeader(req *rpc.Request) error {
	var tmp [1]uint8
	_, err := io.ReadFull(c.RequestReader, tmp[:])
	if err != nil {
		return err
	}

	methodLen := int(tmp[0])
	buff := make([]byte, methodLen)
	_, err = io.ReadAtLeast(c.RequestReader, buff, methodLen)
	if err != nil {
		return err
	}

	req.ServiceMethod = string(buff)
	return nil
}

func (c *serverCodec) ReadRequestBody(v any) error {
	return msgpack.NewDecoder(c.RequestReader).Decode(v)
}

func (c *serverCodec) WriteResponse(resp *rpc.Response, body any) error {
	if resp.Error != "" {
		_, err := c.ResponseWriter.Write([]byte{0xff})
		if err != nil {
			return err
		}
		_, err = io.WriteString(c.ResponseWriter, resp.Error)
		return err
	}

	_, err := c.ResponseWriter.Write([]byte{0xfe})
	if err != nil {
		return err
	}
	return msgpack.NewEncoder(c.ResponseWriter).Encode(body)
}

func (c *serverCodec) Close() error { return nil }

type Server struct {
	rpcServer         *rpc.Server
	serverPrivKey     *bcrypt.NoisePrivateKey
	clientPubKeys     map[string]*bcrypt.NoisePublicKey
	clientPubKeysLock *sync.RWMutex
}

func NewServer() *Server {
	s := &Server{
		rpcServer:         rpc.NewServer(),
		clientPubKeys:     make(map[string]*bcrypt.NoisePublicKey),
		clientPubKeysLock: new(sync.RWMutex),
	}
	return s
}

func (s *Server) SetServerPrivateKey(pk bcrypt.NoisePrivateKey) {
	s.serverPrivKey = &pk
}

func (s *Server) AddClientPublicKey(pk bcrypt.NoisePublicKey) {
	s.clientPubKeysLock.Lock()
	defer s.clientPubKeysLock.Unlock()

	h := string(hash(pk[:]))
	s.clientPubKeys[h] = &pk
}

func (s *Server) RemoveClientPublicKey(pk bcrypt.NoisePublicKey) {
	s.clientPubKeysLock.Lock()
	defer s.clientPubKeysLock.Unlock()

	h := string(hash(pk[:]))
	delete(s.clientPubKeys, h)
}

func (s *Server) Register(rcvr any) error {
	return s.rpcServer.Register(rcvr)
}

func (s *Server) RegisterName(name string, rcvr any) error {
	return s.rpcServer.RegisterName(name, rcvr)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	headerE := req.Header.Get("X-Rpc-E")
	headerT := req.Header.Get("X-Rpc-T")
	if headerE == "" || headerT == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	t, err := strconv.ParseInt(headerT, 10, 64)
	if err != nil || time.Now().Unix()-t > 3*60 {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var ePubKey bcrypt.NoisePublicKey
	if err := ePubKey.FromString(headerE); err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	data, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err = s.process(&ePubKey, t, data, &ePubKey, &t)
	if err != nil {
		if _, ok := err.(internalError); ok {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusForbidden)
		}
		return
	}

	w.Header().Set("X-Rpc-E", ePubKey.String())
	w.Header().Set("X-Rpc-T", strconv.FormatInt(t, 10))
	w.Write(data)
}

func (s *Server) ServeListener(l net.Listener) error {
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go func() {
			defer conn.Close()
			s.ServeConn(conn)
		}()
	}
}

func (s *Server) ServeConn(rw io.ReadWriter) error {
	var (
		ePubKey bcrypt.NoisePublicKey
		t       int64
		dataLen uint16
		data    []byte
	)

	_, err := io.ReadFull(rw, ePubKey[:])
	if err != nil {
		return err
	}

	err = binary.Read(rw, binary.BigEndian, &t)
	if err != nil {
		return err
	}

	if time.Now().Unix()-t > 3*60 {
		return errors.New("request expired")
	}

	err = binary.Read(rw, binary.BigEndian, &dataLen)
	if err != nil {
		return err
	}

	data = make([]byte, dataLen)
	_, err = io.ReadFull(rw, data)
	if err != nil {
		return err
	}

	data, err = s.process(&ePubKey, t, data, &ePubKey, &t)
	if err != nil {
		return err
	}
	dataLen = uint16(len(data))

	if _, err = rw.Write(ePubKey[:]); err != nil {
		return err
	}

	if err = binary.Write(rw, binary.BigEndian, t); err != nil {
		return err
	}

	if err = binary.Write(rw, binary.BigEndian, dataLen); err != nil {
		return err
	}

	_, err = rw.Write(data)
	return err
}

func (s *Server) process(
	ePubKeyIn *bcrypt.NoisePublicKey, tIn int64, dataIn []byte,
	ePubKeyOut *bcrypt.NoisePublicKey, tOut *int64,
) (dataOut []byte, err error) {
	// 解密数据
	data, err := decrypt(s.serverPrivKey, ePubKeyIn, tIn, dataIn)
	if err != nil {
		return nil, err
	}

	// 前8字节是ClientPublicKey的hash
	// 后8字节是ClientPrivateKey * ServerPublicKey的hash
	if len(data) < 16 {
		return nil, io.ErrUnexpectedEOF
	}

	// 查看ClientPublicKey的hash是否在列表里
	s.clientPubKeysLock.RLock()
	clientPubKey, ok := s.clientPubKeys[string(data[:8])]
	s.clientPubKeysLock.RUnlock()
	if !ok {
		return nil, errors.New("client public key invalid")
	}

	// 验证ClientPublicKey是否有效
	ss := s.serverPrivKey.SharedSecret(clientPubKey)
	if !bcrypt.Equals(hash(ss[:]), data[8:16]) {
		return nil, errors.New("client public key invalid")
	}

	// 先创建一个加密用的临时密钥，不然等rpc处理完了才发现有错，就很讨厌
	ePrivKey, err := bcrypt.NewPrivateKey()
	if err != nil {
		return nil, internalError(err)
	}

	// 调用rpc处理，并获取处理后的结果数据
	responseWriter := bytes.NewBuffer(dataIn[:0])
	codec := &serverCodec{
		RequestReader:  bytes.NewReader(data[16:]),
		ResponseWriter: responseWriter,
	}
	err = s.rpcServer.ServeRequest(codec)
	if err != nil {
		return nil, internalError(err)
	}

	// 对结果数据进行加密
	t := time.Now().Unix()
	data = responseWriter.Bytes()
	data, err = encrypt(&ePrivKey, clientPubKey, t, data)
	if err != nil {
		return nil, internalError(err)
	}

	*ePubKeyOut = ePrivKey.PublicKey()
	*tOut = t
	return data, nil
}
