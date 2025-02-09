package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lincketheo/ndbgo/internal/bytecode"
)

func printHelp() {
	fmt.Printf(
		`Usage: %s [SCRIPT]
    If a database name is provided, 
    create or connect to the specified database.
    If no database name is provided, 
    you'll need to connect to the database yourself

Options:
`, os.Args[0])
	flag.PrintDefaults()
}

func main() {
	// Define command line arguments
	help := flag.Bool("help", false, "Show this help message")
	scanout := flag.String("sout", "", "Specify the output file of the scanning process")

  // Define help menu
	flag.Usage = printHelp
	flag.Parse()

  // Print help if present
	if *help {
		printHelp()
		return
	}

  // Get positional arguments
	rest := flag.Args()
	if len(rest) != 1 {
		fmt.Print("Must provide one input file\n")
		return
	}
  data, err := os.ReadFile(rest[0])
  if err != nil {
    fmt.Println(err)
    return
  }

  // Output for scan step
	if len(*scanout) > 0 {
		tokens := bytecode.Scan(string(data))
    if err := bytecode.WriteTokensToFileClean(*scanout, tokens); err != nil {
      fmt.Println(err)
    }
	}
}
