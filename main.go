package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"unsafe"

	"github.com/lincketheo/numstore/internal/numstore"
)

func main() {
	// Create empty connection
	// Create Database
	err := numstore.CreateDatabase("foo")

	if err != nil {
		log.Fatal(err)
	}

	err = numstore.CreateVariable("foo", numstore.VariableMeta{
		Name:  "fiz",
		Dtype: numstore.U32,
		Shape: []uint32{},
	})
	if err != nil {
		log.Fatal(err)
	}

	err = numstore.CreateVariable("foo", numstore.VariableMeta{
		Name:  "bar",
		Dtype: numstore.U64,
		Shape: []uint32{},
	})
	if err != nil {
		log.Fatal(err)
	}

	err = numstore.CreateVariable("foo", numstore.VariableMeta{
		Name:  "baz",
		Dtype: numstore.U32,
		Shape: []uint32{},
	})
	if err != nil {
		log.Fatal(err)
	}

  fiz := []uint32{1, 9, 10, 12}
  bar := []uint64{10, 11, 12, 13}
  baz := []uint32{9, 8, 7, 6}

  buffer := make([]byte, 0, unsafe.Sizeof(fiz) + unsafe.Sizeof(bar) + unsafe.Sizeof(baz))
  for i := range 4 {
    buffer = binary.LittleEndian.AppendUint32(buffer, fiz[i])
    buffer = binary.LittleEndian.AppendUint64(buffer, bar[i])
  }
  for i := range 4 {
    buffer = binary.LittleEndian.AppendUint32(buffer, baz[i])
  }
  
  r := bytes.NewReader(buffer)

  w := numstore.WriteContextRequest{
    ToWrite: 4,
    Dbname: "foo",
    Variables: [][]string{{"fiz", "bar"}, {"baz"}},
  }

  ow, err := w.OpenWriteContextRequest()
  if err != nil {
    log.Fatal(err)
  }
  wow, err := ow.OpenWriteContext()
  if err != nil {
    log.Fatal(err)
  }

  err = wow.WriteAll(r)
  if err != nil {
    log.Fatal(err)
  }

}
