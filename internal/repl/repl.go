package repl

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/lincketheo/ndbgo/internal/usecases"
	"github.com/lincketheo/ndbgo/internal/utils"
)

func replPreambleDb(
	db string,
	n *usecases.NDB,
) error {
	if err := (*n).ConnectDB(db); err != nil {
		return utils.ErrorContext(err)
	}
	return nil
}

func replPreambleRel(
	db,
	rel string,
	n *usecases.NDB,
) error {
	if err := (*n).ConnectDB(db); err != nil {
		return utils.ErrorContext(err)
	} else if err := (*n).ConnectRel(rel); err != nil {
		return utils.ErrorContext(err)
	}
	return nil
}

func replPreambleVari(
	db,
	rel,
	vari string,
	n *usecases.NDB,
) error {
	if err := (*n).ConnectDB(db); err != nil {
		return utils.ErrorContext(err)
	} else if err := (*n).ConnectRel(rel); err != nil {
		return utils.ErrorContext(err)
	} else if err := (*n).ConnectVar(vari); err != nil {
		return utils.ErrorContext(err)
	}
	return nil
}

func preamble(arg string, n *usecases.NDB) error {
	if len(arg) == 0 {
		return nil
	}

	parts := strings.SplitN(arg, ":", 3)

	switch len(parts) {
	case 1:
		db := parts[0]
		if err := replPreambleDb(db, n); err != nil {
			return utils.ErrorContext(err)
		}
		break
	case 2:
		db, rel := parts[0], parts[1]
		if err := replPreambleRel(db, rel, n); err != nil {
			return utils.ErrorContext(err)
		}
		break
	case 3:
		db, rel, vari := parts[0], parts[1], parts[2]
		if err := replPreambleVari(db, rel, vari, n); err != nil {
			return utils.ErrorContext(err)
		}
		break
	default:
		return fmt.Errorf("Invalid argument: %s\n", arg)
	}

	return nil
}

func RunREPL(arg string, n *usecases.NDB) error {
	if err := preamble(arg, n); err != nil {
		return utils.ErrorContext(err)
	}

  fmt.Println("Type 'exit' to quit")
  scanner := bufio.NewScanner(os.Stdin)

	for {
    fmt.Print("> ")
    if !scanner.Scan() {
      break
    }

    input := scanner.Text()
    if input == "exit" {
      fmt.Println("Goodbye!")
      break
    }

    fmt.Println(input)
	}

  if err := scanner.Err(); err != nil {
    return utils.ErrorContext(err)
  }

	return nil
}
