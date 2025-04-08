package compiler

/*
Range string:
- start:stop:step
- start:
- :stop
- start::step

a[0]
a[0, 1, 2] - each index(i) < dim[i]
a[:, 1, 2]
a[:]
a[0:1:2]
*/

type dimRange struct {
	start        int
	stop         int
	step         int
	startPresent bool
	stopPresent  bool
	stepPresent  bool
}

type arrRange struct {
  ranges []dimRange
}
