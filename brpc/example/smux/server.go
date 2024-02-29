package main

import (
	"errors"
	"net"
	"os"
	"time"

	"github.com/abxuz/b-tools/bcrypt"
	"github.com/abxuz/b-tools/brpc"
	"github.com/spf13/cobra"
	"github.com/xtaci/smux"
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

			l, err := net.Listen("tcp", ":10000")
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			for {
				conn, err := l.Accept()
				if err != nil {
					break
				}

				go func() {
					defer conn.Close()
					session, err := smux.Server(conn, nil)
					if err != nil {
						return
					}
					defer session.Close()

					for {
						stream, err := session.AcceptStream()
						if err != nil {
							break
						}
						go func() {
							defer stream.Close()
							stream.SetDeadline(time.Now().Add(time.Second * 3))
							rpcServer.ServeConn(stream)
						}()
					}
				}()
			}
		},
	}
	return c
}
