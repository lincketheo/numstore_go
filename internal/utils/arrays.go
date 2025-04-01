package utils

func ReduceMultU32(arr []uint32) uint32 {
	var ret uint32 = 1
	for _, i := range arr {
		ret *= i
	}
	return ret
}
