package main

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func main() {
	c := &cobra.Command{
		Use: filepath.Base(os.Args[0]),
	}
	c.AddCommand(NewServerCmd())
	c.AddCommand(NewClientCmd())
	c.Execute()
}
