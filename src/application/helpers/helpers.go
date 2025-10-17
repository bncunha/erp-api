package helper

import "strconv"

func ParseInt64(value string) int64 {
	intValue, _ := strconv.ParseInt(value, 10, 64)
	return intValue
}

func ParseFloat(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}
