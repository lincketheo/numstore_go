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
  if err := compiler.RunTokens(compiler.Scan(cmd)); err != nil {
    fmt.Println(err)
  }
}
