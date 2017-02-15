package chyle

import (
	"fmt"
	"regexp"
	"strconv"
)

// RegexpExtracter use a regexp to extract data
type RegexpExtracter struct {
	index      string
	identifier string
	re         *regexp.Regexp
}

// Extract data from a commitMap
func (r RegexpExtracter) Extract(commitMap *map[string]interface{}) (*map[string]interface{}, error) {
	if value, ok := (*commitMap)[r.index]; ok {
		v, ok := value.(string)

		if !ok {
			return nil, fmt.Errorf(`Can't parse value`)
		}

		var result string

		results := r.re.FindStringSubmatch(v)

		if len(results) > 1 {
			result = results[1]
		}

		b, err := strconv.ParseBool(result)

		if err == nil {
			(*commitMap)[r.identifier] = b

			return commitMap, nil
		}

		i, err := strconv.ParseInt(result, 10, 64)

		if err == nil {
			(*commitMap)[r.identifier] = i

			return commitMap, nil
		}

		f, err := strconv.ParseFloat(result, 64)

		if err == nil {
			(*commitMap)[r.identifier] = f

			return commitMap, nil
		}

		(*commitMap)[r.identifier] = result

		return commitMap, nil
	}

	return commitMap, nil
}
