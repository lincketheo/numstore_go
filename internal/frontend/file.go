package nsfrontend

import (
	"log"

	"github.com/lincketheo/numstore/internal/utils"
)

func FileRun(fname, dbname string) {
	contents, err := utils.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
		return
	}
	handleCmd(string(contents))
}
