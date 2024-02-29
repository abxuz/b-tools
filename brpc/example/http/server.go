package main

import (
	"errors"
	"net/http"

	"github.com/abxuz/b-tools/bcrypt"
	"github.com/abxuz/b-tools/brpc"
	"github.com/spf13/cobra"
)

type Service struct {
}

type QueryRequest struct {
	Name string
}

type QueryResponse struct {
	Age int
}

func (s *Service) Query(req QueryRequest, resp *QueryResponse) error {
	if req.Name != "admin" {
		return errors.New("unkown name")
	}
	resp.Age = 100
	return nil
}

func NewServerCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "server",
		Short: "server demo",
		Run: func(cmd *cobra.Command, args []string) {
			var serverPrivKey bcrypt.NoisePrivateKey
			serverPrivKey.FromString("qHTyTvwGYKFeww0tn0/Gdn7vkPvfAsfSUFeXwNUCpnU=")

			var clientPubKey bcrypt.NoisePublicKey
			clientPubKey.FromString("ns1Wlf1dcYaE1gRsgPwU5hy6Kl/psRk6qV84JF24fQI=")

			rpcServer := brpc.NewServer()
			rpcServer.SetServerPrivateKey(serverPrivKey)
			rpcServer.AddClientPublicKey(clientPubKey)
			rpcServer.RegisterName("service", &Service{})

			mux := http.NewServeMux()
			mux.Handle("/rpc", rpcServer)
			http.ListenAndServe(":10001", mux)
		},
	}
	return c
}
