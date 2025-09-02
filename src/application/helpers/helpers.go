package helper

import "strconv"

func ParseInt64(value string) int64 {
	intValue, _ := strconv.ParseInt(value, 10, 64)
	return intValue
}
