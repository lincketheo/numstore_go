package numstore

type Numstore interface {
	HandleCreate(vname string, tp Type) error
	HandleDelete(vname string) error
	HandleRead(fmt ReadFormat) error
	HandleWrite(fmt WriteFormat) error
	HandleTake(fmt ReadFormat) error
}

