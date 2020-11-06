package main

import (
	"fmt"
	"github.com/dmoles/adler/server"
	"github.com/spf13/cobra"
	"os"
)

const defaultPort = 8181

func main() {
	var port int
	var cssDir string
	var cmd = &cobra.Command{
		Use:  "adler <rootDir>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.Start(port, args[0], cssDir)
		},
	}
	cmd.Flags().IntVarP(&port, "port", "p", defaultPort, "port to listen on")
	cmd.Flags().StringVar(&cssDir, "css", "", "CSS/SCSS directory (for testing)")
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
