package function

import (
	"strings"
)

func Trim(str string) (trimmed string, err error) {
	trimmed = strings.ReplaceAll(str, " ", "")
	return
}
