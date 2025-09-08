package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yansol0/blacklight/parser"
	"github.com/yansol0/blacklight/reporter"
	"github.com/yansol0/blacklight/tester"
)

func main() {
	specPath := flag.String("spec", "", "Path to OpenAPI spec JSON")
	baseURL := flag.String("base-url", "", "Base URL of the API")
	token := flag.String("token", "", "Valid JWT token for authenticated probes")
	cookie := flag.String("cookie", "", "Valid session cookie for authenticated probes")
	outdir := flag.String("outdir", "reports", "Directory to write reports")
	flag.Parse()

	if *specPath == "" || *baseURL == "" {
		fmt.Println("Usage: --spec spec.json --base-url https://api.example.com [--token <jwt>] [--cookie <cookie>] --outdir reports")
		os.Exit(1)
	}
	if *token == "" && *cookie == "" {
		fmt.Println("You must provide either --token or --cookie for authenticated probes")
		os.Exit(1)
	}

	endpoints, err := parser.ParseOpenAPISpec(*specPath, *baseURL)
	if err != nil {
		fmt.Println("Failed to parse spec:", err)
		os.Exit(1)
	}

	results := tester.RunTests(endpoints, *token, *cookie)

	err = reporter.WriteReports(results, *outdir)
	if err != nil {
		fmt.Println("Failed to write reports:", err)
		os.Exit(1)
	}
}
