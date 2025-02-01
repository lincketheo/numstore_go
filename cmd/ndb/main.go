package main

import (
	"os"

	"github.com/lincketheo/ndbgo/internal/logging"
	"github.com/lincketheo/ndbgo/internal/repl"
	"github.com/lincketheo/ndbgo/internal/usecases"
)

func main() {
	var _n usecases.NDBlogger
	var n usecases.NDB = &_n

	if len(os.Args) > 2 {
		logging.Error("Usage: %s [DB]:[REL]:[VAR]\n", os.Args[0])
		return
	}

	arg := ""
	if len(os.Args) == 2 {
		arg = os.Args[1]
	}

	if err := repl.RunREPL(arg, &n); err != nil {
		logging.Error("%v", err)
	}
}
