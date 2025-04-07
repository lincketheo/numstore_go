package main

import (
	"os"

	nsfrontend "github.com/lincketheo/numstore/internal/frontend"
)

func main() {
	args := os.Args
	switch len(args) {
	case 2:
		nsfrontend.ShellRun(args[1])
	case 3:
		nsfrontend.FileRun(args[1], args[2])
	default:
		os.Exit(1)
	}
}
