package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/dmoles/adler/server"
)

const defaultPort = 8181

func main() {
	var port int
	var cssDir string
	var cmd = &cobra.Command{
		Use:  "adler <rootDir>",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var dir string
			if len(args) > 0 {
				dir = args[0]
			} else {
				dir = "."
			}
			return server.Start(port, dir, cssDir)
		},
	}
	cmd.Flags().IntVarP(&port, "port", "p", defaultPort, "port to listen on")
	cmd.Flags().StringVar(&cssDir, "css", "", "alternate CSS directory (must contain main.css)")
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
