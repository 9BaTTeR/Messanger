package gencode

import (
	"math/rand"
	"strconv"
)

type Code string

func GenCode() Code {

	c := Code(strconv.FormatUint(uint64(gencode()), 10))
	return c
}
func gencode() uint {
	code := rand.Intn(999999-100000) + 100000
	return uint(code)
}
