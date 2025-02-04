package main

/*
#cgo CFLAGS: -I../../c
#cgo LDFLAGS: -L../../c -lndbc
#include "ndbc.h"
*/
import "C"
import (
	"os"
)

func main() {
  fd := int(os.Stdin.Fd())

	C.write_data(C.int(fd))
}
