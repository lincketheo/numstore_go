package utils

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/lincketheo/ndbgo/internal/config"
)

func IntPlusWillOverflow(a int, b int) bool {
	return a > 0 && b > math.MaxInt-a
}

func IntPlusWillUnderflow(a int, b int) bool {
	return a < 0 && b < math.MinInt-a
}

func IntCanPlus(a int, b int) bool {
	return !IntPlusWillOverflow(a, b) && !IntPlusWillUnderflow(a, b)
}

func CanUint64BeUint32(n uint64) bool {
	return n <= uint64(^uint32(0))
}

func CanIntBeByte(n int) bool {
	return n >= 0 && n <= 255
}

func UInt32ArrBytes(arr []uint32) []byte {
	ret := make([]byte, len(arr)*4)

	for i, v := range arr {
		offset := i * 4
		config.Endian.PutUint32(ret[offset:], uint32(v))
	}

	return ret
}

func BytesToUInt32Arr(b []byte) []uint32 {
	ASSERT(len(b)%4 == 0)

	arr := make([]uint32, len(b)/4)
	for i := range arr {
		offset := i * 4
		arr[i] = config.Endian.Uint32(b[offset:])
	}

	return arr
}

func ParseUInt32Slice(input string) ([]uint32, error) {
	var result []uint32
	err := json.Unmarshal([]byte(input), &result)
	return result, err
}

func IntToUint32(arr []int) ([]uint32, error) {
	ret := make([]uint32, len(arr), len(arr))

	for i, v := range arr {
		if v < 0 {
			return nil, fmt.Errorf("Expected values > 0, got: %d\n", v)
		}
		ret[i] = uint32(v)
	}

	return ret, nil
}
