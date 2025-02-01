package utils

import (
	"github.com/lincketheo/ndbgo/internal/config"
)

func ASSERT(expr bool) {
	if config.Debug {
		if !expr {
			panic("Assertion failed")
		}
	}
}

func UNREACHABLE() {
  panic("Unreachable code reached")
}
