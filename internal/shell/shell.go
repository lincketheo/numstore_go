package shell

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/lincketheo/numstore/internal/compiler"
)

func displayHelp() {
	fmt.Println(".help    - Show available commands")
	fmt.Println(".clear   - Clear the terminal screen")
	fmt.Println(".exit    - Closes your connection")
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func cleanInput(text string) string {
	output := strings.TrimSpace(text)
	return output
}

func handleCmd(cmd string) {
  compiler.Parse(compiler.Scan(cmd))
}

func Run(dbName string) {
	controlCommands := map[string]any{
		".help":  displayHelp,
		".clear": clearScreen,
	}

	reader := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for reader.Scan() {
		text := cleanInput(reader.Text())
		if command, exists := controlCommands[text]; exists {
			command.(func())()
		} else if strings.EqualFold(".exit", text) {
			return
		} else {
			handleCmd(text)
		}
    fmt.Print("> ")
	}
  fmt.Println()
}
