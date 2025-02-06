package utils

import "runtime"

func getFunctionName(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return "Unknown"
	}
	return runtime.FuncForPC(pc).Name()
}
