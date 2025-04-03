package utils

import "runtime"

func GetFunctionName(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return "Unknown"
	}
	return runtime.FuncForPC(pc).Name()
}
