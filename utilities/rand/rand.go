package rand

import (
	"math/rand"
	"time"
)

var (
	randSource rand.Source
)

func init() {
	rand.Seed(time.Now().UnixNano())
	randSource = rand.NewSource(time.Now().UnixNano())
}

func Intn(maxValue int) int {
	return rand.Intn(maxValue)
}
