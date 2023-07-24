package helper

import (
	"fmt"
	"strings"
)

func PadLeft(origin string, overall int, char string) (result string) {
	result = origin
	if overall < len(origin) {
		return
	}
	deviation := overall - len(origin)
	result = fmt.Sprintf("%s%s", strings.Repeat(char, deviation), origin)
	return
}

func PadRight(origin string, overall int, char string) (result string) {
	result = origin
	if overall < len(origin) {
		return
	}
	deviation := overall - len(origin)
	result = fmt.Sprintf("%s%s", origin, strings.Repeat(char, deviation))
	return
}
