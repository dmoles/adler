package main

import (
	"fmt"
	"github.com/dmoles/adler/adler"
	"github.com/spf13/cobra"
	"os"
)

const defaultPort = 8080

func startCmd() *cobra.Command {
	var port int
	cmd := &cobra.Command{
		Use: "start <root-dir>",
		Short: "Start Adler server",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return start(port, args[0])
		},
	}
	cmd.Flags().IntVarP(&port, "port", "p", defaultPort, "port to listen on")
	return cmd
}

func start(port int, rootDir string) error {
	server, err := adler.NewServer(port, rootDir)
	if err != nil {
		return err
	}
	return server.Start()
}

func main() {
	var rootCmd = &cobra.Command{
		Use: "adler",
	}
	rootCmd.AddCommand(startCmd())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
