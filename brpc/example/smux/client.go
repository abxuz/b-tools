package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/abxuz/b-tools/bcrypt"
	"github.com/abxuz/b-tools/brpc/rw"
	"github.com/spf13/cobra"
	"github.com/xtaci/smux"
)

func NewClientCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "client",
		Short: "client demo",
		Run: func(cmd *cobra.Command, args []string) {
			var clientPrivKey bcrypt.NoisePrivateKey
			clientPrivKey.FromString("iBv818sWwMDjU/IdvVyb2hAvlTrm6S/xf9oSFySEVnw=")

			var serverPubKey bcrypt.NoisePublicKey
			serverPubKey.FromString("7S7lkXbp3Xomf9WdCbvL68hxEcdGxT4X+Wco4gKa2CM=")

			conn, err := net.Dial("tcp", "127.0.0.1:10000")
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			defer conn.Close()

			session, err := smux.Client(conn, nil)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}
			defer session.Close()

			rpcClient := rw.NewClient(
				rw.WithClientPrivateKey(clientPrivKey),
				rw.WithServerPublicKey(serverPubKey),
				rw.WithOpen(func() (io.ReadWriteCloser, error) {
					stream, err := session.OpenStream()
					if err != nil {
						return nil, err
					}
					stream.SetDeadline(time.Now().Add(time.Second * 3))
					return stream, nil
				}),
			)

			var resp *QueryResponse
			err = rpcClient.Call("service.Query", QueryRequest{Name: "123"}, &resp)
			if err != nil {
				fmt.Println(err)
			} else {
				cmd.PrintErrln("unexpected result")
				os.Exit(1)
			}

			err = rpcClient.Call("service.Query", QueryRequest{Name: "admin"}, &resp)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			fmt.Println(resp)
		},
	}
	return c
}
