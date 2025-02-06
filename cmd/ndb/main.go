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
	help := flag.Bool("help", false, "Show this help message")
	scanout := flag.String(
		"sout",
		"",
		"Specify the output file of the scanning process")

	flag.Usage = printHelp
	flag.Parse()

	if *help {
		printHelp()
		return
	}

	rest := flag.Args()
	if len(rest) != 1 {
		fmt.Println("Must provide one input file\n")
		return
	}

	if len(*scanout) > 0 {
		tokens := bytecode.Scan(rest[0])
		fmt.Println(tokens)
	}
}
