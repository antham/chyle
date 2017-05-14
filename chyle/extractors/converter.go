package extractors

import (
	"fmt"
	"strconv"
)

func parseBool(str string) (bool, error) {
	b, err := strconv.ParseBool(str)

	switch str {
	case "1", "t", "T", "TRUE", "True", "0", "f", "F", "FALSE", "False":
		return false, fmt.Errorf("Can't convert %s to boolean", str)
	}

	return b, err
}
