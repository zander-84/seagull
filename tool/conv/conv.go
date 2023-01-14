package conv

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"strconv"
)

func ShouldStrToInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return n
}

func ShouldStrToInt32(s string) int32 {
	int10, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0
	}
	return int32(int10)
}

func ShouldStringToInt64(s string) int64 {
	int10, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return int10
}

func ShouldStrToUint(u string) uint {
	uInt64, err := strconv.ParseUint(u, 10, 0)
	if err != nil {
		return 0
	}
	return uint(uInt64)
}

func ShouldStrToUint32(u string) uint32 {
	uInt64, err := strconv.ParseUint(u, 10, 32)
	if err != nil {
		return 0
	}
	return uint32(uInt64)
}

func ShouldStrToF64(s string) float64 {
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return n
}

func F64toStr(f float64) string {
	return fmt.Sprintf("%f", f)
}

func IntegerToStr[T constraints.Integer](in T) string {
	return strconv.FormatInt(int64(in), 10)
}

func SliceIntegerToSliceStr[T constraints.Integer](in []T) []string {
	var idsStr = make([]string, 0, len(in))
	for _, id := range in {
		idsStr = append(idsStr, IntegerToStr(id))
	}
	return idsStr
}

// StrToBool
// 字符串bool转系统bool
// 备注：
// True类型:"1", "t", "T", "true", "TRUE", "True"
// False类型:"0", "f", "F", "false", "FALSE", "False"
func StrToBool(b string) bool {
	boolVal, err := strconv.ParseBool(b)
	if err != nil {
		return false
	}
	return boolVal
}

func ShouldF64ToDecimal(f float64, len int) float64 {
	format := fmt.Sprintf("%%.%df", len)
	value, err := strconv.ParseFloat(fmt.Sprintf(format, f), 64)
	if err != nil {
		return 0
	}
	return value
}
