package chyle

import (
	"fmt"
	"regexp"

	"github.com/antham/envh"
)

// extracter describe a way to extract data from a commit hashmap summary
type extracter interface {
	extract(*map[string]interface{}) (*map[string]interface{}, error)
}

// extract parse commit fields to extract datas
func extract(extractors *[]extracter, commitMaps *[]map[string]interface{}) (*Changelog, error) {
	var err error

	results := []map[string]interface{}{}

	for _, commitMap := range *commitMaps {
		result := &commitMap

		for _, extractor := range *extractors {
			result, err = extractor.extract(result)

			if err != nil {
				return nil, err
			}
		}

		results = append(results, *result)
	}

	return &Changelog{Datas: results, Metadatas: map[string]interface{}{}}, nil
}

// createExtractors build extracters from a config
func createExtractors(config *envh.EnvTree) (*[]extracter, error) {
	results := []extracter{}

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
			return &[]extracter{}, fmt.Errorf(`"%s" is not a valid regular expression defined for "EXTRACTORS_%s_%s" key`, datas["REG"], identifier, "REG")
		}

		debug(`Extractor "%s" "ORIGKEY" defined with value "%s"`, identifier, datas["ORIGKEY"])
		debug(`Extractor "%s" "DESTKEY" defined with value "%s"`, identifier, datas["DESTKEY"])
		debug(`Extractor "%s" "REG" defined with value "%s"`, identifier, datas["REG"])

		results = append(results, regexpExtractor{
			datas["ORIGKEY"],
			datas["DESTKEY"],
			re,
		})
	}

	return &results, nil
}
