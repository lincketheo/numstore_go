package main

import (
	"github.com/lincketheo/ndbgo/internal/logging"
	"github.com/lincketheo/ndbgo/internal/repl"
)

func main() {
  if err := repl.RunREPL("foo:bar"); err != nil {
		logging.Error("%v", err)
	}
}
