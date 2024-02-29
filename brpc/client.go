package brpc

import (
	"bytes"
	"context"
	"errors"
	"io"
	"time"

	"github.com/abxuz/b-tools/bcrypt"
	"github.com/vmihailenco/msgpack/v5"
)

type LogicError string

func (e LogicError) Error() string { return string(e) }

type RpcClient interface {
	Call(serviceName string, req any, resp any) error
	CallContext(ctx context.Context, serviceName string, req any, resp any) error
}

type Client struct {
	clientPrivKey    *bcrypt.NoisePrivateKey
	serverPubKey     *bcrypt.NoisePublicKey
	clientPubKeyHash []byte
	ssHash           []byte
}

func (c *Client) SetClientPrivateKey(pk bcrypt.NoisePrivateKey) {
	clientPubKey := pk.PublicKey()
	c.clientPrivKey = &pk
	c.clientPubKeyHash = hash(clientPubKey[:])
	c.tryFixSsHash()
}

func (c *Client) SetServerPublicKey(pk bcrypt.NoisePublicKey) {
	c.serverPubKey = &pk
	c.tryFixSsHash()
}

func (c *Client) WriteRequestMessage(
	dst []byte, serviceName string, req any,
	ePubKeyOut *bcrypt.NoisePublicKey, tOut *int64,
) (dataOut []byte, err error) {
	ePrivKey, err := bcrypt.NewPrivateKey()
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(dst)
	buffer.Write(c.clientPubKeyHash)
	buffer.Write(c.ssHash)
	buffer.WriteByte(byte(len(serviceName)))
	buffer.WriteString(serviceName)
	err = msgpack.NewEncoder(buffer).Encode(req)
	if err != nil {
		return nil, err
	}

	data := buffer.Bytes()
	t := time.Now().Unix()
	data, err = encrypt(&ePrivKey, c.serverPubKey, t, data)
	if err != nil {
		return nil, err
	}

	*ePubKeyOut = ePrivKey.PublicKey()
	*tOut = t
	return data, nil
}

func (c *Client) ReadResponseMessage(resp any, data []byte, ePubKey *bcrypt.NoisePublicKey, t int64) error {
	data, err := decrypt(c.clientPrivKey, ePubKey, t, data)
	if err != nil {
		return err
	}

	if len(data) < 1 {
		return io.ErrUnexpectedEOF
	}

	if data[0] == 0xff {
		return LogicError(data[1:])
	}

	if data[0] != 0xfe {
		return errors.New("invalid response data")
	}

	return msgpack.Unmarshal(data[1:], resp)
}

func (c *Client) tryFixSsHash() {
	if c.serverPubKey == nil || c.clientPrivKey == nil {
		return
	}

	ss := c.clientPrivKey.SharedSecret(c.serverPubKey)
	c.ssHash = hash(ss[:])
}
