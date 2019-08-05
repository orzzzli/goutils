package convert

import "strconv"

func Int64to32(number int64) (int,error) {
	tempStr := strconv.FormatInt(number, 10)
	return strconv.Atoi(tempStr)
}

func Int32to64(number int) (int64,error) {
	tempStr := strconv.Itoa(number)
	return strconv.ParseInt(tempStr, 10, 64)
}