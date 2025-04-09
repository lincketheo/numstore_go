package nsfrontend

import (
	"fmt"
	"os"
	"strings"

	"github.com/lincketheo/numstore/internal/compiler"
)

func cleanInput(text string) string {
	output := strings.TrimSpace(text)
	return output
}

func handleCmd(cmd string) {
  if bytes, ok := compiler.Parse(compiler.Scan(cmd)); !ok {
    fmt.Println("ERROR")
  } else {
    os.Stdout.Write(bytes)
  }
}
