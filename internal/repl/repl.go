package repl

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
	"github.com/lincketheo/ndbgo/internal/dtypes"
	"github.com/lincketheo/ndbgo/internal/logging"
	"github.com/lincketheo/ndbgo/internal/usecases"
	"github.com/lincketheo/ndbgo/internal/utils"
)

func RunREPL(arg string, n *usecases.NDB) error {
	var err error

	if arg != "" {
		if err := replConnect([]string{arg}, n); err != nil {
			return utils.ErrorContext(err)
		}
	}

	fmt.Println("Type 'exit' to quit")

	rl, err := readline.NewEx(&readline.Config{
		Prompt:      "> ",
		HistoryFile: "/tmp/repl_history.txt",
	})

	if err != nil {
		return utils.ErrorContext(err)
	}

	defer func() {
		if cerr := rl.Close(); cerr != nil {
			err = cerr
		}
	}()

	for {
		input, err := rl.Readline()
		if err != nil {
			break
		}

		if input == "exit" {
			fmt.Println("Goodbye!")
			break
		}

		if err := replCommand(input, n); err != nil {
			logging.Error("\n%v", utils.ErrorContext(err))
		}
	}

	return err
}

func shellCommand(args []string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func replCommand(cmd string, n *usecases.NDB) error {
	args := strings.Fields(cmd)

	if len(args) < 1 {
		return nil
	}

	arg := strings.ToLower(args[0])
	switch arg {
	case "create":
		return replCreate(args[1:], n)
	case "connect":
		return replConnect(args[1:], n)
	default:
		return shellCommand(args) 
	}
}

// ////////////////////////////// CONNECT

func replConnect(args []string, n *usecases.NDB) error {
	if len(args) == 0 {
		return nil
	}

	name, rest := args[0], args[1:]

	if len(rest) > 0 {
		return fmt.Errorf("Expected no arguments for Connect, but got: %v", rest)
	}

	p, err := parseDbRelVarStr(name)
	if err != nil {
		return utils.ErrorContext(err)
	}

	switch p.present {
	case 1:
		if err := replConnectDb(p.db, n); err != nil {
			return utils.ErrorContext(err)
		}
		break
	case 2:
		if err := replConnectRel(p.db, p.rel, n); err != nil {
			return utils.ErrorContext(err)
		}
		break
	case 3:
		if err := replConnectVar(p.db, p.rel, p.vari, n); err != nil {
			return utils.ErrorContext(err)
		}
		break
	default:
		panic("Unreachable!")
	}

	return nil
}

func replConnectDb(
	db string,
	n *usecases.NDB,
) error {
	if err := (*n).ConnectDB(db); err != nil {
		return utils.ErrorContext(err)
	}
	return nil
}

func replConnectRel(
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

func replConnectVar(
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

// ////////////////////////////// CREATE

func replCreate(args []string, n *usecases.NDB) error {
	if len(args) == 0 {
		return fmt.Errorf("Expected argument after create")
	}

	name, rest := args[0], args[1:]

	parts, err := parseDbRelVarStr(name)
	if err != nil {
		return utils.ErrorContext(err)
	}

	switch parts.present {
	case 1:
		if err := replCreateDb(
			parts.db,
			rest, n); err != nil {
			return utils.ErrorContext(err)
		}
		break
	case 2:
		if err := replCreateRel(
			parts.db,
			parts.rel,
			rest, n); err != nil {
			return utils.ErrorContext(err)
		}
		break
	case 3:
		if err := replCreateVar(
			parts.db,
			parts.rel,
			parts.vari,
			rest, n); err != nil {
			return utils.ErrorContext(err)
		}
		break
	default:
		panic("Unreachable!")
	}

	return nil
}

func replCreateDb(
	db string,
	rest []string,
	n *usecases.NDB,
) error {
	if len(rest) != 0 {
		return fmt.Errorf("Expected no arguments for Create DB, but got: %v", rest)
	}

	if err := (*n).CreateDB(db); err != nil {
		return utils.ErrorContext(err)
	}
	return nil
}

func replCreateRel(
	db,
	rel string,
	rest []string,
	n *usecases.NDB,
) error {
	if len(rest) != 0 {
		return fmt.Errorf("Expected no arguments for Create REL, but got: %v", rest)
	}

	if err := (*n).ConnectDB(db); err != nil {
		return utils.ErrorContext(err)
	} else if err := (*n).CreateRel(rel); err != nil {
		return utils.ErrorContext(err)
	}
	return nil
}

func replCreateVar(
	db,
	rel,
	vari string,
	rest []string,
	n *usecases.NDB,
) error {
	if config, err := replParseVarConfig(rest); err != nil {
		return err
	} else {
		return usecases.CreateDBRelVar(db, rel, vari, *config, n)
	}
}

// ////////////////////////////// UTILS

func parseUint32Args(args []string) ([]uint32, error) {
	if len(args) == 0 {
		return []uint32{}, nil
	}

	ret := make([]uint32, len(args))

	for i, sizestr := range args {
		if size, err := strconv.ParseUint(sizestr, 10, 32); err != nil {
			return nil, utils.ErrorContext(err)
		} else {
			utils.ASSERT(utils.CanUint64BeUint32(size))
			ret[i] = uint32(size)
		}
	}

	return ret, nil
}

func replParseVarConfig(config []string) (*usecases.VarConfig, error) {
	if len(config) == 0 {
		return nil, fmt.Errorf("Expected DTYPE")
	} else if dtype, err := dtypes.StrtoDtype(config[0]); err != nil {
		return nil, utils.ErrorContext(err)
	} else if shape, err := parseUint32Args(config[1:]); err != nil {
		return nil, utils.ErrorContext(err)
	} else {
		return &usecases.VarConfig{
			Dtype: dtype,
			Shape: shape,
		}, nil
	}
}

type dbRelVar struct {
	db      string
	rel     string
	vari    string
	present int
}

func parseDbRelVarStr(name string) (*dbRelVar, error) {
	parts := strings.SplitN(name, ":", 3)

	if len(parts) > 0 && len(parts[0]) == 0 {
		return nil, fmt.Errorf("DB Name must not be empty")
	}
	if len(parts) > 1 && len(parts[1]) == 0 {
		return nil, fmt.Errorf("Rel Name must not be empty")
	}
	if len(parts) > 2 && len(parts[2]) == 0 {
		return nil, fmt.Errorf("Var Name must not be empty")
	}

	switch len(parts) {
	case 1:
		return &dbRelVar{
			db:      parts[0],
			rel:     "",
			vari:    "",
			present: 1,
		}, nil
	case 2:
		return &dbRelVar{
			db:      parts[0],
			rel:     parts[1],
			vari:    "",
			present: 2,
		}, nil
	case 3:
		return &dbRelVar{
			db:      parts[0],
			rel:     parts[1],
			vari:    parts[2],
			present: 3,
		}, nil
	default:
		return nil, fmt.Errorf("Invalid entity string: %s\n", name)
	}
}
