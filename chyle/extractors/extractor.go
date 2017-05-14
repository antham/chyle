package extractors

import (
	"regexp"

	"github.com/antham/chyle/chyle/types"
)

// Extracter describe a way to extract data from a commit hashmap summary
type Extracter interface {
	Extract(*map[string]interface{}) *map[string]interface{}
}

// Extract parses commit fields to extract datas
func Extract(extractors *[]Extracter, commitMaps *[]map[string]interface{}) *types.Changelog {
	results := []map[string]interface{}{}

	for _, commitMap := range *commitMaps {
		result := &commitMap

		for _, extractor := range *extractors {
			result = extractor.Extract(result)
		}

		results = append(results, *result)
	}

	changelog := types.Changelog{}
	changelog.Datas = results
	changelog.Metadatas = map[string]interface{}{}

	return &changelog
}

// Create builds extracters from a config
func Create(extractors map[string]struct {
	ORIGKEY string
	DESTKEY string
	REG     *regexp.Regexp
}) *[]Extracter {
	results := []Extracter{}

	for _, extractor := range extractors {
		results = append(results, regexpExtractor{
			extractor.ORIGKEY,
			extractor.DESTKEY,
			extractor.REG,
		})
	}

	return &results
}
