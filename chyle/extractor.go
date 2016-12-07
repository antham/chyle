package chyle

import (
	"fmt"
	"regexp"
	"strconv"
)

// Extracter describe a way to extract data from a commit hashmap summary
type Extracter interface {
	Extract(*map[string]interface{}) (*map[string]interface{}, error)
}

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

// Extract parse commit fields to extract datas
func Extract(extractors *[]Extracter, commitMaps *[]map[string]interface{}) (*[]map[string]interface{}, error) {
	var err error

	results := []map[string]interface{}{}

	for _, commitMap := range *commitMaps {
		result := &commitMap

		for _, extractor := range *extractors {
			result, err = extractor.Extract(result)

			if err != nil {
				return nil, err
			}
		}

		results = append(results, *result)
	}

	return &results, nil
}

// CreateExtractors build extracters from a config
func CreateExtractors(extracters map[string]interface{}) (*[]Extracter, error) {
	results := []Extracter{}

	for dk, dv := range extracters {
		e, ok := dv.(map[string]interface{})

		if !ok {
			return &[]Extracter{}, fmt.Errorf(`extractor "%s" must contains key=value string values`, dk)
		}

		for key, value := range e {
			s, ok := value.(string)

			if !ok {
				return &[]Extracter{}, fmt.Errorf(`extractor "%s" is not a string`, s)
			}

			re, err := regexp.Compile(s)

			if err != nil {
				return &[]Extracter{}, fmt.Errorf(`"%s" doesn't contain a valid regular expression`, key)
			}

			results = append(results, RegexpExtracter{
				dk,
				key,
				re,
			})
		}

	}

	return &results, nil
}
