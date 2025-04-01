package main

import (
	"log"

	"github.com/lincketheo/numstore/internal/numstore"
)

func main() {
	// Create empty connection
	con := numstore.EmptyConnection()

	// Create Database
	err := numstore.CreateDatabase(con, numstore.CreateDatabaseArgs{
		Name: "foo",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Connect to database
	con, err = con.ConnectDB("foo")
	if err != nil {
		log.Fatal(err)
	}

	err = numstore.CreateVariable(con, numstore.CreateVariableArgs{
		Name:  "fiz",
		Dtype: numstore.U32,
		Shape: []uint32{1, 2},
	})
	if err != nil {
		log.Fatal(err)
	}
}
