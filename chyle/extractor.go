package chyle

import (
	"fmt"
	"regexp"

	"github.com/antham/envh"
)

// Extracter describe a way to extract data from a commit hashmap summary
type Extracter interface {
	Extract(*map[string]interface{}) (*map[string]interface{}, error)
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
func CreateExtractors(config *envh.EnvTree) (*[]Extracter, error) {
	results := []Extracter{}

	for _, identifier := range config.GetChildrenKeys() {
		subConfig, err := config.FindSubTree(identifier)

		if err != nil {
			return &results, err
		}

		datas := map[string]string{}

		for _, v := range []string{"ORIGKEY", "DESTKEY", "REG"} {
			datas[v], err = subConfig.FindString(v)

			if err != nil {
				return &results, fmt.Errorf(`An environment variable suffixed with "%s" must be defined with "%s", like EXTRACTORS_%s_%s`, v, identifier, identifier, v)
			}
		}

		re, err := regexp.Compile(datas["REG"])

		if err != nil {
			return &[]Extracter{}, fmt.Errorf(`"%s" is not a valid regular expression defined for "EXTRACTORS_%s_%s" key`, datas["REG"], identifier, "REG")
		}

		debug(`Extractor "%s" "ORIGKEY" defined with value "%s"`, identifier, datas["ORIGKEY"])
		debug(`Extractor "%s" "DESTKEY" defined with value "%s"`, identifier, datas["DESTKEY"])
		debug(`Extractor "%s" "REG" defined with value "%s"`, identifier, datas["REG"])

		results = append(results, RegexpExtracter{
			datas["ORIGKEY"],
			datas["DESTKEY"],
			re,
		})
	}

	return &results, nil
}
