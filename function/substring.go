package function

import "errors"

func Substring(str string, args ...int) (subbed string, err error) {
	if len(args) == 1 {
		subbed = str[args[0]:]
	} else if len(args) == 2 {
		subbed = str[args[0]:args[1]]
	} else {
		err = errors.New("No match argument")
	}
	return
}
