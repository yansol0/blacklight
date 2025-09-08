package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yansol0/blacklight/parser"
	"github.com/yansol0/blacklight/reporter"
	"github.com/yansol0/blacklight/tester"
	"github.com/yansol0/blacklight/tui"
	"github.com/yansol0/blacklight/utils"
)

func main() {
	specPath := flag.String("spec", "", "Path to OpenAPI spec JSON")
	baseURLFlag := flag.String("base-url", "", "Base URL of the API (optional, overrides spec)")
	token := flag.String("token", "", "Valid JWT token for authenticated probes")
	cookie := flag.String("cookie", "", "Valid session cookie for authenticated probes")
	outdir := flag.String("outdir", "reports", "Directory to write reports")
	noTUI := flag.Bool("no-tui", false, "Disable interactive TUI (enabled by default)")
	flag.Parse()

	if *specPath == "" {
		fmt.Println("Usage: --spec spec.json [--base-url https://api.example.com] [--token <jwt>] [--cookie <cookie>] --outdir reports [--no-tui]")
		os.Exit(1)
	}
	if *token == "" && *cookie == "" {
		fmt.Println("You must provide either --token or --cookie for authenticated probes")
		os.Exit(1)
	}

	utils.PrintBanner()

	endpoints, resolvedBaseURL, err := parser.ParseOpenAPISpec(*specPath, *baseURLFlag)
	if err != nil {
		fmt.Println("Failed to parse spec:", err)
		os.Exit(1)
	}
	fmt.Println("Using base URL:", resolvedBaseURL)

	var results tester.Results
	if !*noTUI {
		utils.SetLoggingEnabled(false)
		updates := make(chan tui.ProgressUpdate)
		done := make(chan struct{})

		go func() {
			results = tester.RunTestsWithProgress(endpoints, *token, *cookie, updates)
			close(updates)
			close(done)
		}()

		if err := tui.Run(updates, done); err != nil {
			fmt.Println("TUI error:", err)
			os.Exit(1)
		}
	} else {
		results = tester.RunTests(endpoints, *token, *cookie)
	}

	if err := reporter.WriteReports(results, *outdir); err != nil {
		fmt.Println("Failed to write reports:", err)
		os.Exit(1)
	}
}
