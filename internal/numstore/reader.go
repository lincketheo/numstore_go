package numstore

type VariableSlice struct {
	vname string
}

type ReadContextRequest struct {
	ToRead         int
	Dbname         string
	VariableSlices [][]string
}
