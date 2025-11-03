package rw

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"io"
	"time"

	"github.com/abxuz/b-tools/v2/bcrypt"
	"github.com/abxuz/b-tools/v2/brpc"
)

type OpenFunc = func() (io.ReadWriteCloser, error)

type Client struct {
	brpc.Client
	open OpenFunc
}

type option = func(c *Client)

func NewClient(opts ...option) *Client {
	c := new(Client)
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func WithOpen(open OpenFunc) option {
	return func(c *Client) {
		c.open = open
	}
}

func WithServerPublicKey(pk bcrypt.NoisePublicKey) option {
	return func(c *Client) {
		c.SetServerPublicKey(pk)
	}
}

func WithClientPrivateKey(pk bcrypt.NoisePrivateKey) option {
	return func(c *Client) {
		c.SetClientPrivateKey(pk)
	}
}

func (c *Client) Call(serviceName string, req any, resp any) error {
	return c.CallContext(context.Background(), serviceName, req, resp)
}

func (c *Client) CallContext(ctx context.Context, serviceName string, req any, resp any) error {
	var (
		ePubKey bcrypt.NoisePublicKey
		t       int64
	)

	data, err := c.Client.WriteRequestMessage(nil, serviceName, req, &ePubKey, &t)
	if err != nil {
		return err
	}
	dataLen := uint16(len(data))

	rwc, err := c.open()
	if err != nil {
		return err
	}
	defer rwc.Close()

	if _, err := rwc.Write(ePubKey[:]); err != nil {
		return err
	}

	if err := binary.Write(rwc, binary.BigEndian, t); err != nil {
		return err
	}

	if err := binary.Write(rwc, binary.BigEndian, dataLen); err != nil {
		return err
	}

	if _, err := rwc.Write(data); err != nil {
		return err
	}

	if _, err := io.ReadFull(rwc, ePubKey[:]); err != nil {
		return err
	}

	if err := binary.Read(rwc, binary.BigEndian, &t); err != nil {
		return err
	}

	if time.Now().Unix()-t > 3*60 {
		return errors.New("response expired, sync time with server")
	}

	if err := binary.Read(rwc, binary.BigEndian, &dataLen); err != nil {
		return err
	}

	buffer := bytes.NewBuffer(data[:0])
	if _, err := io.CopyN(buffer, rwc, int64(dataLen)); err != nil {
		return err
	}

	data = buffer.Bytes()
	return c.Client.ReadResponseMessage(resp, data, &ePubKey, t)
}
