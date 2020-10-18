package main

import (
	"fmt"
	"github.com/dmoles/adler/server"
	"github.com/spf13/cobra"
	"os"
)

const defaultPort = 8181

func start(port int, rootDir string) error {
	return server.Start(port, rootDir)
}

func main() {
	var port int
	var cmd = &cobra.Command{
		Use:  "adler <rootDir>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return start(port, args[0])
		},
	}
	cmd.Flags().IntVarP(&port, "port", "p", defaultPort, "port to listen on")
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
