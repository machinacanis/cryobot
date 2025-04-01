package utils

import "math/rand"

func RandomDeviceNumber() int {
	return rand.Intn(9999999-1000000+1) + 1000000
}
