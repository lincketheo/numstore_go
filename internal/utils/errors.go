package utils

import (
	"fmt"

	"github.com/lincketheo/ndbgo/internal/config"
)

func ErrorContext(err error) error {
	if err == nil {
		return nil
	}

	if config.Debug {
		caller := getFunctionName(2)
		return fmt.Errorf("%s\n%w", caller, err)
	} else {
		return err
	}
}

func ErrorMoref(cause error, msg string, args ...any) error {
	ASSERT(cause != nil)
	msgRet := fmt.Sprintf(msg, args)
	return fmt.Errorf(msgRet+": %w", cause)
}
