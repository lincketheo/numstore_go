package utils

const enableAsserts = true

func Assert(cond bool) {
	if enableAsserts {
		if !cond {
			panic("Assertion failed!")
		}
	}
}
