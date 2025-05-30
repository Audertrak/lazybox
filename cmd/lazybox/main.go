package main

import (
	"fmt"
	"lazybox/internal/code"
	"lazybox/internal/enuminfo"
	"lazybox/internal/env"
	"lazybox/internal/file"
	"lazybox/internal/fn"
	"lazybox/internal/fs"
	"lazybox/internal/listinfo"
	"lazybox/internal/output"
	"lazybox/internal/structinfo"
	"lazybox/internal/text"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	printBanner()

	var flagAll bool
	var flagVerbose bool
	var flagLess bool
	var flagMin bool
	var flagIncremental bool
	var flagIR bool
	var flagSilent bool
	var flagTokenize bool

	var rootCmd = &cobra.Command{
		Use:   "lazybox",
		Short: "lazybox - swiss army knife for data extraction and formatting",
	}

	rootCmd.PersistentFlags().BoolVarP(&flagAll, "all", "a", false, "Print all representations of the data, including all available metadata and results")
	rootCmd.PersistentFlags().BoolVarP(&flagVerbose, "verbose", "v", false, "Verbose output. Includes additional metadata and results.")
	rootCmd.PersistentFlags().BoolVarP(&flagLess, "less", "l", false, "Compact, minimal output, with selective exclusions of metadata or results.")
	rootCmd.PersistentFlags().BoolVarP(&flagMin, "min", "m", false, "Remove all whitespace and convert to a single string value.")
	rootCmd.PersistentFlags().BoolVarP(&flagIncremental, "incremental", "i", false, "Print the output incrementally as it is processed.")
	rootCmd.PersistentFlags().BoolVarP(&flagIR, "ir", "I", false, "Print the intermediate representation of the data.")
	rootCmd.PersistentFlags().BoolVarP(&flagSilent, "silent", "s", false, "Create an intermediate representation of the data, but do not print it to stdout.")
	rootCmd.PersistentFlags().BoolVarP(&flagTokenize, "tokenize", "t", false, "Remove articles or other prose grammar and use simple key:value pairs.")

	var fsCmd = &cobra.Command{
		Use:   "fs [path] [mode]",
		Short: "Emit a representation of the filesystem given a path",
		Args:  cobra.RangeArgs(0, 2),
		Run: func(cmd *cobra.Command, args []string) {
			path := "."
			if len(args) > 0 {
				path = args[0]
			}
			mode := "jsonify"
			if len(args) > 1 {
				mode = args[1]
			}
			fileInfoIR, err := fs.Scan(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error scanning %s: %v\n", path, err)
				os.Exit(1)
			}
			handleOutput(fileInfoIR, mode, flagIR, flagAll, flagMin, flagLess, flagTokenize, flagSilent)
		},
	}

	var fileCmd = &cobra.Command{
		Use:   "file [path] [mode]",
		Short: "Open and read the contents of a file",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			mode := "jsonify"
			if len(args) > 1 {
				mode = args[1]
			}
			fileDataIR, err := file.Read(path) // Assuming file.Read returns a type compatible with handleOutput (e.g. *ir.FileInfo or *ir.TextInfo)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", path, err)
				os.Exit(1)
			}
			handleOutput(fileDataIR, mode, flagIR, flagAll, flagMin, flagLess, flagTokenize, flagSilent)
		},
	}

	var apiCmd = &cobra.Command{
		Use:   "api [path] [mode]",
		Short: "Extract an API from a source file",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			// path := args[0]
			// mode := "jsonify"
			// if len(args) > 1 {
			// 	mode = args[1]
			// }
			// apiInfoIR, err := api.Extract(path)
			// if err != nil {
			// 	fmt.Fprintf(os.Stderr, "Error extracting API from %s: %v\n", path, err)
			// 	os.Exit(1)
			// }
			// handleOutput(apiInfoIR, mode, flagIR, flagAll, flagMin, flagLess, flagTokenize, flagSilent)
			fmt.Println("TODO: API command output handling to be implemented with generalized handleOutput")
		},
	}

	var pkgCmd = &cobra.Command{
		Use:   "pkg [path] [mode]",
		Short: "Crawl a package directory and emit a representation of its structure",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			// path := args[0]
			// mode := "jsonify"
			// if len(args) > 1 {
			// 	mode = args[1]
			// }
			// pkgInfoIR, err := pkg.Crawl(path)
			// if err != nil {
			// 	fmt.Fprintf(os.Stderr, "Error crawling package %s: %v\n", path, err)
			// 	os.Exit(1)
			// }
			// handleOutput(pkgInfoIR, mode, flagIR, flagAll, flagMin, flagLess, flagTokenize, flagSilent)
			fmt.Println("TODO: PKG command output handling to be implemented with generalized handleOutput")
		},
	}

	var textCmd = &cobra.Command{
		Use:   "text [path] [mode]",
		Short: "Parse a text file and extract its contents and metadata.",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			mode := "jsonify"
			if len(args) > 1 {
				mode = args[1]
			}

			// First, read the file to get its basic info and content
			fileData, err := file.Read(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", path, err)
				os.Exit(1)
			}
			if fileData.Error != "" {
				fmt.Fprintf(os.Stderr, "Error reading file content %s: %v\n", path, fileData.Error)
				os.Exit(1)
			}

			// Now, analyze the content to get TextInfo
			textInfoIR, err := text.Analyze(fileData.Content, path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error analyzing text from %s: %v\n", path, err)
				os.Exit(1)
			}
			// text.Analyze populates its own FileInfo, so we can pass textInfoIR directly.
			handleOutput(textInfoIR, mode, flagIR, flagAll, flagMin, flagLess, flagTokenize, flagSilent)
		},
	}

	var codeCmd = &cobra.Command{
		Use:   "code [path] [mode]",
		Short: "Parse source code and extract relevant information.",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			mode := "jsonify"
			if len(args) > 1 {
				mode = args[1]
			}
			ir, err := code.ExtractPlaceholder(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			handleOutput(ir, mode, flagIR, flagAll, flagMin, flagLess, flagTokenize, flagSilent)
		},
	}

	var funcCmd = &cobra.Command{
		Use:   "func [path] [mode]",
		Short: "Parse a function and extract its signature and parameters.",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			mode := "jsonify"
			if len(args) > 1 {
				mode = args[1]
			}
			ir, err := fn.ExtractPlaceholder(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			handleOutput(ir, mode, flagIR, flagAll, flagMin, flagLess, flagTokenize, flagSilent)
		},
	}

	var envCmd = &cobra.Command{
		Use:   "env [mode]",
		Short: "Parse environment variables and emit their values.",
		Args:  cobra.RangeArgs(0, 1),
		Run: func(cmd *cobra.Command, args []string) {
			mode := "jsonify"
			if len(args) > 0 {
				mode = args[0]
			}
			ir, err := env.ExtractPlaceholder()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			handleOutput(ir, mode, flagIR, flagAll, flagMin, flagLess, flagTokenize, flagSilent)
		},
	}

	var structCmd = &cobra.Command{
		Use:   "struct [path] [mode]",
		Short: "Parse a data structure and emit its contents.",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			mode := "jsonify"
			if len(args) > 1 {
				mode = args[1]
			}
			ir, err := structinfo.ExtractPlaceholder(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			handleOutput(ir, mode, flagIR, flagAll, flagMin, flagLess, flagTokenize, flagSilent)
		},
	}

	var enumCmd = &cobra.Command{
		Use:   "enum [path] [mode]",
		Short: "Parse an enumeration and emit its values.",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			mode := "jsonify"
			if len(args) > 1 {
				mode = args[1]
			}
			ir, err := enuminfo.ExtractPlaceholder(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			handleOutput(ir, mode, flagIR, flagAll, flagMin, flagLess, flagTokenize, flagSilent)
		},
	}

	var listCmd = &cobra.Command{
		Use:   "list [path] [mode]",
		Short: "Parse a data structure and emit its compile-time contents.",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			mode := "jsonify"
			if len(args) > 1 {
				mode = args[1]
			}
			ir, err := listinfo.ExtractPlaceholder(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			handleOutput(ir, mode, flagIR, flagAll, flagMin, flagLess, flagTokenize, flagSilent)
		},
	}

	var fetchCmd = &cobra.Command{
		Use:   "fetch",
		Short: "Show system info in fastfetch/neofetch style",
		Run: func(cmd *cobra.Command, args []string) {
			output.PrintFastfetchStyle()
		},
	}

	rootCmd.AddCommand(fsCmd)
	rootCmd.AddCommand(fileCmd)
	rootCmd.AddCommand(apiCmd)
	rootCmd.AddCommand(pkgCmd)
	rootCmd.AddCommand(textCmd)
	rootCmd.AddCommand(codeCmd)
	rootCmd.AddCommand(funcCmd)
	rootCmd.AddCommand(envCmd)
	rootCmd.AddCommand(structCmd)
	rootCmd.AddCommand(enumCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(fetchCmd)
	rootCmd.Execute()
}

func handleOutput(data interface{}, mode string, flagIR, flagAll, flagMin, flagLess, flagTokenize, flagSilent bool) {
	if data == nil {
		fmt.Fprintln(os.Stderr, "Error: No data to output.")
		os.Exit(1)
	}
	if flagSilent {
		return // Do not print anything
	}

	// IR flag takes precedence and prints the raw data structure passed to handleOutput.
	if flagIR {
		output.PrintJSON(data)
		return
	}

	// The rest of the flags and modes will operate on the `data interface{}`.
	// Output functions (PrintJSON, PrintMarkdown, etc.) must be updated to handle `interface{}`
	// using type assertions or reflection to correctly process different IR types.

	if flagAll {
		output.PrintJSON(data)
		output.PrintMarkdown(data)
		output.PrintPretty(data)
		output.PrintTable(data)
		return
	}

	if flagMin {
		output.PrintMinJSON(data)
		return
	}

	if flagLess {
		output.PrintLessJSON(data)
		return
	}

	if flagTokenize {
		output.PrintTokenized(data)
		return
	}

	switch mode {
	case "jq":
		output.PrintJQ(data)
	case "prettify":
		output.PrintPretty(data)
	// Add other modes here, ensuring they can handle `interface{}`
	default:
		output.Print(data, mode) // Assuming output.Print will handle the interface{} correctly.
	}
}

// Add a CLI banner at startup for supreme vibes
func printBanner() {
	output.PrintBannerText("lazybox")
}
