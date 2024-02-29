package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/abxuz/b-tools/bcrypt"
	"github.com/abxuz/b-tools/brpc"
)

type Client struct {
	brpc.Client
	httpClient *http.Client
	endpoint   string
}

type option = func(c *Client)

func NewClient(opts ...option) *Client {
	c := new(Client)
	for _, opt := range opts {
		opt(c)
	}

	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}

	return c
}

func WithHttpClient(httpClient *http.Client) option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func WithEndpoint(endpoint string) option {
	return func(c *Client) {
		c.endpoint = endpoint
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

	buffer := bytes.NewBuffer(data)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, buffer)
	if err != nil {
		return err
	}

	request.Header.Set("X-Rpc-E", ePubKey.String())
	request.Header.Set("X-Rpc-T", strconv.FormatInt(t, 10))

	response, err := c.httpClient.Do(request)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %v, %v", response.StatusCode, response.Status)
	}

	headerE := response.Header.Get("X-Rpc-E")
	headerT := response.Header.Get("X-Rpc-T")
	if headerE == "" || headerT == "" {
		return errors.New("invalid response, required header missing")
	}

	if err := ePubKey.FromString(headerE); err != nil {
		return err
	}

	t, err = strconv.ParseInt(headerT, 10, 64)
	if err != nil {
		return err
	}
	if time.Now().Unix()-t > 3*60 {
		return errors.New("response expired, sync time with server")
	}

	buffer.Reset()
	_, err = io.Copy(buffer, response.Body)
	if err != nil {
		return err
	}

	data = buffer.Bytes()
	return c.Client.ReadResponseMessage(resp, data, &ePubKey, t)
}
