package nsfrontend

import (
	"fmt"
	"strings"

	"github.com/lincketheo/numstore/internal/compiler"
)

func cleanInput(text string) string {
	output := strings.TrimSpace(text)
	return output
}

func handleCmd(cmd string) {
  if !compiler.Parse(compiler.Scan(cmd)) {
    fmt.Println("ERROR")
  }
}
