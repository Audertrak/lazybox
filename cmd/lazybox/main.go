package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"lazybox/internal/fs"
	"lazybox/internal/output"
	"os"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "lazybox",
		Short: "lazybox - swiss army knife for data extraction and formatting",
	}

	var fsCmd = &cobra.Command{
		Use:   "fs [path]",
		Short: "Emit a representation of the filesystem given a path",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			path := "."
			if len(args) > 0 {
				path = args[0]
			}
			ir, err := fs.Scan(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			output.PrintJSON(ir)
		},
	}

	rootCmd.AddCommand(fsCmd)
	rootCmd.Execute()
}
