package main

import (
	"fmt"
	"github.com/dmoles/adler/adler"
	"github.com/spf13/cobra"
	"os"
)

const defaultPort = 8181

func start(port int, rootDir string) error {
	server, err := adler.NewServer(port, rootDir)
	if err != nil {
		return err
	}
	return server.Start()
}

func main() {
	var port int
	var cmd = &cobra.Command{
		Use: "adler",
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
