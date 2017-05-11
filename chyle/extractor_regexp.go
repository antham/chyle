package chyle

import (
	"regexp"
	"strconv"
)

// regexpExtractor use a regexp to extract data
type regexpExtractor struct {
	index      string
	identifier string
	re         *regexp.Regexp
}

// extract data from a commitMap
func (r regexpExtractor) extract(commitMap *map[string]interface{}) *map[string]interface{} {
	var mapValue interface{}
	var ok bool

	if mapValue, ok = (*commitMap)[r.index]; !ok {
		return commitMap
	}

	var value string

	value, ok = mapValue.(string)

	if !ok {
		return commitMap
	}

	var result string

	results := r.re.FindStringSubmatch(value)

	if len(results) > 1 {
		result = results[1]
	}

	b, err := parseBool(result)

	if err == nil {
		(*commitMap)[r.identifier] = b

		return commitMap
	}

	i, err := strconv.ParseInt(result, 10, 64)

	if err == nil {
		(*commitMap)[r.identifier] = i

		return commitMap
	}

	f, err := strconv.ParseFloat(result, 64)

	if err == nil {
		(*commitMap)[r.identifier] = f

		return commitMap
	}

	(*commitMap)[r.identifier] = result

	return commitMap
}
