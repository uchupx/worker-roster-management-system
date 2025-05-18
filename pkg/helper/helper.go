package helper

import "strconv"

func StringToInt(str string) int64 {
	parsedInt, err := strconv.Atoi(str)
	if err != nil {
		return 0 
	}
	return int64(parsedInt)
}
