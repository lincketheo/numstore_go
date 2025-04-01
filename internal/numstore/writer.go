package numstore

import "github.com/lincketheo/numstore/internal/nserror"

type WriteContextRequest struct {
	variables [][]string
}

type WriteContext struct {
	variables [][]Variable
}

type OpenWriteContext struct {
	variables [][]WritingVariable
}

func (w WriteContextRequest) toWriteContext(con Connection) (WriteContext, error) {
	result := make([][]Variable, len(w.variables))
	usedVars := make(map[string]struct{})

	for i, vcol := range w.variables {
		varCol := make([]Variable, len(vcol))
		for j, v := range vcol {

			// Do not allow duplicate variables
			if _, exists := usedVars[v]; exists {
				return WriteContext{}, nserror.WriteContext_DuplicateVar
			}

			// Load variable meta
			variable, err := LoadVariable(con, v)
			if err != nil {
				return WriteContext{}, err
			}

			varCol[j] = variable
		}
		result[i] = varCol
	}

	return WriteContext{variables: result}, nil
}
func (w WriteContext) Open(con Connection) (OpenWriteContext, error) {
	result := make([][]WritingVariable, len(w.variables))

	for i, vcol := range w.variables {
		varCol := make([]WritingVariable, len(vcol))
		for j, v := range vcol {

			// Load and open variable
			openVar, err := v.Open(con)
			if err != nil {
				goto failed
			}

			varCol[j] = openVar
		}
		result[i] = varCol
	}

	return OpenWriteContext{variables: result}, nil

failed:
	// Close all previously opened WritingVariables
  var zero WritingVariable
	for x := range result {
		for y := range result[x] {
			if result[x][y] != zero {
				result[x][y].Close()
			}
		}
	}
	return OpenWriteContext{}, fmt.Errorf("failed to open write context")
}
