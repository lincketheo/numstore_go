package nserror

import (
	"fmt"
	"strings"
)

type DomainError struct {
	Code     int
	Messages []string
	Cause    error
}

func (e *DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s\n%s", strings.Join(e.Messages, "\n"), e.Cause.Error())
	}
	return fmt.Sprintf("[%d] %s", e.Code, strings.Join(e.Messages, "\n"))
}

func New(code int, msg string) error {
	return &DomainError{
		Code:     code,
		Messages: []string{msg},
		Cause:    nil,
	}
}

func extractCode(err error) int {
	if de, ok := err.(*DomainError); ok {
		return de.Code
	}
	return 0
}

func Wrapf(err error, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return &DomainError{
		Code:     extractCode(err),
		Messages: []string{msg},
		Cause:    err,
	}
}

var (
	// Create DB errors
	DBAlreadyExists          = New(1, "database already exists")
	DBNotConnected           = New(2, "database not connected")
	ConnectionInvalidState   = New(3, "invalid connection state")
	DBAlreadyConnected       = New(4, "database already connected")
	VarAlreadyConnected      = New(5, "variable already connected")
	ConnectingToVarWithoutDB = New(6, "cannot connect variable without a database")
	InvalidVarName           = New(7, "invalid variable name")
	InvalidDBName            = New(8, "invalid database name")
	NotConnectedToDB         = New(9, "not connected to database")
	VarAlreadyExists         = New(10, "variable already exists")
	JSONInvalidDtype         = New(11, "invalid dtype for JSON")

	// WriteContext
	WriteContext_DuplicateVar = New(12, "Invalid Write context, duplicate variable")
)
