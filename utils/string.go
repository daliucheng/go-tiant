package utils

import (
	"fmt"
	"strings"
)

func JoinArgs(showByte int, args ...interface{}) string {
	var sumLen int

	argStr := make([]string, len(args))
	for i, v := range args {
		if s, ok := v.(string); ok {
			argStr[i] = s
		} else {
			argStr[i] = fmt.Sprintf("%v", v)
		}

		sumLen += len(argStr[i])
		if sumLen >= showByte {
			break
		}
	}

	argVal := strings.Join(argStr, " ")
	if sumLen > showByte {
		argVal = argVal[:showByte] + " ..."
	}
	return argVal
}
