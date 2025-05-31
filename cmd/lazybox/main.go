// cmd/lazybox/main.go
package main

import (
	"fmt"
	"lazybox/internal/code"
	"lazybox/internal/enuminfo"
	"lazybox/internal/env"
	"lazybox/internal/file"
	"lazybox/internal/fn"
	"lazybox/internal/fs"
	"lazybox/internal/glpg" // Added GLPG import
	"lazybox/internal/listinfo"
	"lazybox/internal/output"
	"lazybox/internal/structinfo"
	"lazybox/internal/text"
	"lazybox/internal/theme" // Import the theme package
	"os"

	"github.com/spf13/cobra"
)

func main() {
	theme.Initialize() // Initialize the theme system
	printBanner()

	var flagAll bool
	var flagVerbose bool
	var flagLess bool
	var flagMin bool
	var flagIncremental bool
	var flagIR bool
	var flagSilent bool
	var flagTokenize bool

	var outputMode string // Variable to hold the output mode from the flag

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
	rootCmd.PersistentFlags().StringVarP(&outputMode, "output", "o", "jsonify", "Output mode (e.g., jsonify, prettify, mdify, tableify, commafy, fastfetch)")

	var fsCmd = &cobra.Command{
		Use:   "fs [path] [mode]",
		Short: "Emit a representation of the filesystem given a path",
		Args:  cobra.RangeArgs(0, 2),
		Run: func(cmd *cobra.Command, args []string) {
			path := "."
			if len(args) > 0 {
				path = args[0]
			}
			// Mode is now determined by the --output flag or its default value "jsonify"
			// The positional mode argument is effectively ignored if --output is used.
			// For now, the --output flag takes precedence.
			mode := outputMode
			if len(args) > 1 {
				// If a positional mode is given AND --output was not changed from default,
				// we can consider using the positional one.
				// However, to keep it simple, we'll let the --output flag control it.
				// If --output is explicitly set, it overrides the positional argument.
				// If --output is not set, it defaults to "jsonify", and if a positional mode is also present,
				// the current logic will use the --output default.
				// To prioritize positional:
				if cmd.Flags().Changed("output") {
					mode = outputMode
				} else {
					mode = args[1]
				}
			}

			fileInfoIR, err := fs.Scan(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error scanning %s: %v\n", path, err)
				os.Exit(1)
			}
			glpgData, err := glpg.ToGLPG(fileInfoIR)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error converting to GLPG: %v\n", err)
				os.Exit(1)
			}
			handleOutput(glpgData, mode, collectFlags(cmd))
		},
	}

	var fileCmd = &cobra.Command{
		Use:   "file [path] [mode]",
		Short: "Open and read the contents of a file",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			mode := outputMode // Use the --output flag
			if len(args) > 1 && !cmd.Flags().Changed("output") {
				mode = args[1] // Fallback to positional if --output not explicitly set
			}
			fileDataIR, err := file.Read(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", path, err)
				os.Exit(1)
			}
			glpgData, err := glpg.ToGLPG(fileDataIR)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error converting to GLPG: %v\n", err)
				os.Exit(1)
			}
			handleOutput(glpgData, mode, collectFlags(cmd))
		},
	}

	var apiCmd = &cobra.Command{
		Use:   "api [path] [mode]",
		Short: "Extract an API from a source file",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			// path := args[0]
			// mode := outputMode // Use the --output flag
			// if len(args) > 1 && !cmd.Flags().Changed("output") {
			// 	mode = args[1] // Fallback to positional
			// }
			// apiInfoIR, err := api.Extract(path) // Placeholder for actual API extraction
			// if err != nil {
			// 	fmt.Fprintf(os.Stderr, "Error extracting API from %s: %v\n", path, err)
			// 	os.Exit(1)
			// }
			// glpgData, err := glpg.ToGLPG(apiInfoIR)
			// if err != nil {
			// 	fmt.Fprintf(os.Stderr, "Error converting API info to GLPG: %v\n", err)
			// 	os.Exit(1)
			// }
			// handleOutput(glpgData, mode, collectFlags(cmd))
			fmt.Println("TODO: API command output handling to be implemented with GLPG")
		},
	}

	var pkgCmd = &cobra.Command{
		Use:   "pkg [path] [mode]",
		Short: "Crawl a package directory and emit a representation of its structure",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			// path := args[0]
			// mode := outputMode // Use the --output flag
			// if len(args) > 1 && !cmd.Flags().Changed("output") {
			// 	mode = args[1] // Fallback to positional
			// }
			// pkgInfoIR, err := pkg.Crawl(path) // Placeholder for actual package crawling
			// if err != nil {
			// 	fmt.Fprintf(os.Stderr, "Error crawling package %s: %v\n", path, err)
			// 	os.Exit(1)
			// }
			// glpgData, err := glpg.ToGLPG(pkgInfoIR)
			// if err != nil {
			// 	fmt.Fprintf(os.Stderr, "Error converting package info to GLPG: %v\n", err)
			// 	os.Exit(1)
			// }
			// handleOutput(glpgData, mode, collectFlags(cmd))
			fmt.Println("TODO: PKG command output handling to be implemented with GLPG")
		},
	}

	var textCmd = &cobra.Command{
		Use:   "text [path] [mode]",
		Short: "Parse a text file and extract its contents and metadata.",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			mode := outputMode // Use the --output flag
			if len(args) > 1 && !cmd.Flags().Changed("output") {
				mode = args[1] // Fallback to positional
			}

			fileData, err := file.Read(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", path, err)
				os.Exit(1)
			}
			if fileData.Error != "" {
				fmt.Fprintf(os.Stderr, "Error reading file content %s: %v\n", path, fileData.Error)
				os.Exit(1)
			}

			textInfoIR, err := text.Analyze(*fileData.Content, path) // Dereference fileData.Content
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error analyzing text from %s: %v\n", path, err)
				os.Exit(1)
			}
			glpgData, err := glpg.ToGLPG(textInfoIR)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error converting to GLPG: %v\n", err)
				os.Exit(1)
			}
			handleOutput(glpgData, mode, collectFlags(cmd))
		},
	}

	var codeCmd = &cobra.Command{
		Use:   "code [path] [mode]",
		Short: "Parse source code and extract relevant information.",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			mode := outputMode // Use the --output flag
			if len(args) > 1 && !cmd.Flags().Changed("output") {
				mode = args[1] // Fallback to positional
			}
			codeInfoIR, err := code.ExtractPlaceholder(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			glpgData, err := glpg.ToGLPG(codeInfoIR)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error converting to GLPG: %v\n", err)
				os.Exit(1)
			}
			handleOutput(glpgData, mode, collectFlags(cmd))
		},
	}

	var funcCmd = &cobra.Command{
		Use:   "func [path] [mode]",
		Short: "Parse a function and extract its signature and parameters.",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			mode := outputMode // Use the --output flag
			if len(args) > 1 && !cmd.Flags().Changed("output") {
				mode = args[1] // Fallback to positional
			}
			funcInfoIR, err := fn.ExtractPlaceholder(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			glpgData, err := glpg.ToGLPG(funcInfoIR)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error converting to GLPG: %v\n", err)
				os.Exit(1)
			}
			handleOutput(glpgData, mode, collectFlags(cmd))
		},
	}

	var envCmd = &cobra.Command{
		Use:   "env [mode]",
		Short: "Parse environment variables and emit their values.",
		Args:  cobra.RangeArgs(0, 1),
		Run: func(cmd *cobra.Command, args []string) {
			mode := outputMode // Use the --output flag
			if len(args) > 0 && !cmd.Flags().Changed("output") {
				mode = args[0] // Fallback to positional
			}
			envInfoIR, err := env.ExtractPlaceholder()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			glpgData, err := glpg.ToGLPG(envInfoIR)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error converting to GLPG: %v\n", err)
				os.Exit(1)
			}
			handleOutput(glpgData, mode, collectFlags(cmd))
		},
	}

	var structCmd = &cobra.Command{
		Use:   "struct [path] [mode]",
		Short: "Parse a data structure and emit its contents.",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			mode := outputMode // Use the --output flag
			if len(args) > 1 && !cmd.Flags().Changed("output") {
				mode = args[1] // Fallback to positional
			}
			structInfoIR, err := structinfo.ExtractPlaceholder(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			glpgData, err := glpg.ToGLPG(structInfoIR)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error converting to GLPG: %v\n", err)
				os.Exit(1)
			}
			handleOutput(glpgData, mode, collectFlags(cmd))
		},
	}

	var enumCmd = &cobra.Command{
		Use:   "enum [path] [mode]",
		Short: "Parse an enumeration and emit its values.",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			mode := outputMode // Use the --output flag
			if len(args) > 1 && !cmd.Flags().Changed("output") {
				mode = args[1] // Fallback to positional
			}
			enumInfoIR, err := enuminfo.ExtractPlaceholder(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			glpgData, err := glpg.ToGLPG(enumInfoIR)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error converting to GLPG: %v\n", err)
				os.Exit(1)
			}
			handleOutput(glpgData, mode, collectFlags(cmd))
		},
	}

	var listCmd = &cobra.Command{
		Use:   "list [path] [mode]",
		Short: "Parse a data structure and emit its compile-time contents.",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			mode := outputMode // Use the --output flag
			if len(args) > 1 && !cmd.Flags().Changed("output") {
				mode = args[1] // Fallback to positional
			}
			listInfoIR, err := listinfo.ExtractPlaceholder(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			glpgData, err := glpg.ToGLPG(listInfoIR)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error converting to GLPG: %v\n", err)
				os.Exit(1)
			}
			handleOutput(glpgData, mode, collectFlags(cmd))
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
	// rootCmd.AddCommand(fetchCmd) // Commented out as fetchCmd is not defined in the provided code
	rootCmd.Execute()
}

// collectFlags gathers all persistent flags into a map.
func collectFlags(cmd *cobra.Command) map[string]bool {
	flags := make(map[string]bool)
	flags["all"], _ = cmd.Flags().GetBool("all")
	flags["verbose"], _ = cmd.Flags().GetBool("verbose")
	flags["less"], _ = cmd.Flags().GetBool("less")
	flags["min"], _ = cmd.Flags().GetBool("min")
	flags["incremental"], _ = cmd.Flags().GetBool("incremental")
	flags["ir"], _ = cmd.Flags().GetBool("ir")
	flags["silent"], _ = cmd.Flags().GetBool("silent")
	flags["tokenize"], _ = cmd.Flags().GetBool("tokenize")
	return flags
}

func handleOutput(data *glpg.GLPG, mode string, flags map[string]bool) {
	if data == nil {
		fmt.Fprintln(os.Stderr, "Error: No data to output.")
		os.Exit(1)
	}
	if silentFlag, ok := flags["silent"]; ok && silentFlag {
		return // Do not print anything
	}

	var err error
	// Dispatch to the correct GLPG-based output function based on mode
	switch mode {
	case "jsonify":
		err = output.PrintGLPGAsJSON(data, flags)
	case "prettify":
		err = output.PrintGLPGAsPretty(data, flags)
	case "mdify":
		err = output.PrintGLPGAsMarkdown(data, flags)
	case "tabelify":
		err = output.PrintGLPGAsTable(data, flags)
	case "commafy": // Added commafy (CSV)
		err = output.PrintGLPGAsCSV(data, flags)
	case "fastfetch": // Added fastfetch style
		err = output.PrintGLPGAsFastfetch(data, flags)
	// Add other cases as they are implemented
	default:
		fmt.Fprintf(os.Stderr, "Warning: Mode '%s' not implemented or recognized for GLPG. Defaulting to JSON output.\n", mode)
		err = output.PrintGLPGAsJSON(data, flags) // Fallback to JSON
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during output generation for mode '%s': %v\n", mode, err)
		// Decide if we should os.Exit(1) here. For now, just printing the error.
	}
}

// Add a CLI banner at startup for supreme vibes
func printBanner() {
	output.PrintBannerText("lazybox")
}
