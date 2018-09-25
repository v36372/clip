package param

import "strconv"

func GetIntFromStrWithDefault(str string, def int) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		return def
	}
	return num
}
