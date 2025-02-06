package main

import (
	"fmt"

	"github.com/lincketheo/ndbgo/internal/bytecode"
)

func main() {
	cmd := `
    create variable []() 12.123 # foobar biz 123 I 12
    foobar 12 I f
  `

	fmt.Println(cmd)

	bytecode.Compile(cmd)
}
