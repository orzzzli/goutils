package util

import "strconv"

func Int64to32(number int64) (int,error) {
	tempStr := strconv.FormatInt(number, 10)
	return strconv.Atoi(tempStr)
}