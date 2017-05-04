package chyle

import (
	"regexp"
)

// extracter describe a way to extract data from a commit hashmap summary
type extracter interface {
	extract(*map[string]interface{}) *map[string]interface{}
}

// extract parse commit fields to extract datas
func extract(extractors *[]extracter, commitMaps *[]map[string]interface{}) *Changelog {
	results := []map[string]interface{}{}

	for _, commitMap := range *commitMaps {
		result := &commitMap

		for _, extractor := range *extractors {
			result = extractor.extract(result)
		}

		results = append(results, *result)
	}

	return &Changelog{Datas: results, Metadatas: map[string]interface{}{}}
}

// createExtractors build extracters from a config
func createExtractors() *[]extracter {
	results := []extracter{}

	for _, datas := range chyleConfig.EXTRACTORS {
		results = append(results, regexpExtractor{
			datas["ORIGKEY"],
			datas["DESTKEY"],
			regexp.MustCompile(datas["REG"]),
		})
	}

	return &results
}
