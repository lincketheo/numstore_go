package numstore

type DimRange struct {
	Start    int
	Stop     int
	Step     int
	IsInf    bool
	IsSingle bool
}

type ReadVariable struct {
	Vname string
	Range []DimRange
}

type ReadFormat struct {
	ToRead    int
	Variables [][]ReadVariable
}
