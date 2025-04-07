package compiler

import "fmt"

func expectedAfter(exp, after tokenType) error {
	return fmt.Errorf("Expected token type %v after token type: %v",
		exp, after)
}

func invalidAfter(after, before tokenType) error {
	return fmt.Errorf("Invalid token type: %v after token type: %v",
		after, before)
}

func invalidAfterExpected(after, before, exp tokenType) error {
	return fmt.Errorf("Invalid token type: %v after token type: %v. "+
		"Expected token type: %v", after, before, exp)
}

func invalidExpected(tok, exp tokenType) error {
	return fmt.Errorf("Invalid token type: %v "+
		"Expected token type: %v", tok, exp)
}

func expectedTokButEOF(tok tokenType) error {
	return fmt.Errorf("Unexpected EOF, expected token type", tok)
}

func expectedStrButEOF(str string) error {
	return fmt.Errorf("Unexpected EOF, expected %s", str)
}

func expectedAnyButGot(got tokenType, exp ...tokenType) error {
	return fmt.Errorf("Expected any token: %v, "+
		"but got token: %v", exp, got)
}
