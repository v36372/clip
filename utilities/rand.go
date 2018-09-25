package utils

import (
	"clip/utilities/rand"
)

func GetRandomNumber(ran int) int {
	return rand.Intn(ran)
}
