package numstore

import (
	"fmt"
	"io"

	"github.com/lincketheo/numstore/internal/nserror"
)

type WriteFormat struct {
	ToWrite   int
	Variables [][]string
}

type WriteFormat2 struct {
	ToWrite   int
	Dbname    string
	Variables [][]string
}

type WriteContext struct {
	toWrite   int
	dbname    string
	variables [][]VariableMeta
}

type OpenWriteContext struct {
	toWrite   int
	dbname    string
	variables [][]WritingVariable
}

func (w WriteFormat2) OpenWriteFormat2() (WriteContext, error) {
	result := make([][]VariableMeta, len(w.Variables))
	usedVars := make(map[string]struct{})

	for i, vcol := range w.Variables {
		varCol := make([]VariableMeta, len(vcol))
		for j, v := range vcol {

			// Do not allow duplicate variables
			if _, exists := usedVars[v]; exists {
        return WriteContext{}, fmt.Errorf("Variable: %s was used more than once in write format string", v)
			}

			// Load variable meta
			variable, err := LoadVariableMeta(w.Dbname, v)
			if err != nil {
				return WriteContext{}, nserror.ErrorStack(err)
			}

			varCol[j] = variable
		}
		result[i] = varCol
	}

	return WriteContext{
		variables: result,
		dbname:    w.Dbname,
		toWrite:   w.ToWrite,
	}, nil
}

func (w WriteContext) OpenWriteContext() (OpenWriteContext, error) {
	result := make([][]WritingVariable, len(w.variables))

	i := 0
	j := 0
	var vcol []VariableMeta
	var v VariableMeta
	var err error
	var openVar WritingVariable

	for i, vcol = range w.variables {
		varCol := make([]WritingVariable, len(vcol))
		for j, v = range vcol {

			// Load and open variable
			openVar, err = OpenVariable(w.dbname, v)
			if err != nil {
				goto failed
			}

			varCol[j] = openVar
		}
		result[i] = varCol
	}

	return OpenWriteContext{
		variables:  result,
		dbname:     w.dbname,
		toWrite:    w.toWrite,
	}, nil

failed:
	for ; i >= 0; i-- {
		for ; j >= 0; j-- {
			result[i][j].Close()
		}
	}
	return OpenWriteContext{}, nserror.ErrorStack(err)
}

func (w OpenWriteContext) WriteAll(r io.Reader) error {
	for _, vcol := range w.variables {

		for range w.toWrite {

			for _, v := range vcol {
				if err := v.WriteNext(r, 0); err != nil {
					return nserror.ErrorStack(err)
				}
			}
		}
	}

	return nil
}

func (w OpenWriteContext) CloseWriteContext() error {
	var err error = nil

	for _, wcol := range w.variables {
		for _, ow := range wcol {
			if _err := ow.Close(); _err != nil {
				err = _err
			}
		}
	}

	return nserror.ErrorStack(err)
}
