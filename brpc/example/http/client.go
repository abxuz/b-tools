package main

import (
	"fmt"
	"os"

	"github.com/abxuz/b-tools/bcrypt"
	"github.com/abxuz/b-tools/brpc/http"
	"github.com/spf13/cobra"
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

			rpcClient := http.NewClient(
				http.WithClientPrivateKey(clientPrivKey),
				http.WithServerPublicKey(serverPubKey),
				http.WithEndpoint("http://127.0.0.1:10001/rpc"),
			)

			var resp *QueryResponse
			err := rpcClient.Call("service.Query", QueryRequest{Name: "123"}, &resp)
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
