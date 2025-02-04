package main 

/*
#cgo CFLAGS: -I../../c
#cgo LDFLAGS: -L../../c -lndbc
#include "ndbc.h"
*/
import "C"

func main() {
  C.read_data() 
}
