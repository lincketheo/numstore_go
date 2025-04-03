package nserror

import (
	"fmt"

	"github.com/lincketheo/numstore/internal/utils"
)

func ErrorContext(err error) error {
	if err == nil {
		return nil
	}

	if utils.Debug {
		caller := utils.GetFunctionName(2)
		return fmt.Errorf("%s\n%w", caller, err)
	} else {
		return err
	}
}

func ErrorMoref(cause error, msg string, args ...any) error {
	utils.Assert(cause != nil)
	msgRet := fmt.Sprintf(msg, args...)
	return fmt.Errorf(msgRet+": %w", cause)
}

func ErrorStack(cause error) error {
  return ErrorMoref(cause, utils.GetFunctionName(2))
}
