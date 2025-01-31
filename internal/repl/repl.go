package repl

import (
	"fmt"
	"strings"

	"github.com/lincketheo/ndbgo/internal/compiler"
	"github.com/lincketheo/ndbgo/internal/interpreter"
	"github.com/lincketheo/ndbgo/internal/ndb"
)

func replPreambleDb(
	db string,
	n *interpreter.NDB,
	p *compiler.Program,
) error {
	if err := p.ConnectDB(db); err != nil {
		return err
	} else if err = interpreter.RunAll(n, p); err != nil {
		return err
	}
	return nil
}

func replPreambleRel(
	db,
	rel string,
	n *interpreter.NDB,
	p *compiler.Program,
) error {
	if err := p.ConnectDB(db); err != nil {
		return err
	} else if err := p.ConnectREL(rel); err != nil {
		return err
	} else if err = interpreter.RunAll(n, p); err != nil {
		return err
	}
	return nil
}

func replPreambleVari(
	db,
	rel,
	vari string,
	n *interpreter.NDB,
	p *compiler.Program,
) error {
	if err := p.ConnectDB(db); err != nil {
		return err
	} else if err := p.ConnectREL(rel); err != nil {
		return err
	} else if err := p.ConnectVar(vari); err != nil {
		return err
	} else if err = interpreter.RunAll(n, p); err != nil {
		return err
	}
	return nil
}

func RunREPL(arg string) error {
	program := compiler.CreateProgram()
	impl := ndb.CreateNDBimpl()
	var interpreter interpreter.NDB = &impl

	parts := strings.SplitN(arg, ":", 3)

	switch len(parts) {
	case 1:
		db := parts[0]
		if err := replPreambleDb(db, &interpreter, &program); err != nil {
			return err
		}
		break
	case 2:
		db, rel := parts[0], parts[1]
		if err := replPreambleRel(db, rel, &interpreter, &program); err != nil {
			return err
		}
		break
	case 3:
		db, rel, vari := parts[0], parts[1], parts[2]
		if err := replPreambleVari(db, rel, vari, &interpreter, &program); err != nil {
			return err
		}
		break
	default:
		return fmt.Errorf("Invalid argument: %s\n", arg)
	}

	return nil
}
