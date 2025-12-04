package extractors

import (
	"github.com/antham/chyle/chyle/types"
)

// Extracter describes a way to extract data from a commit hashmap summary
type Extracter interface {
	Extract(*map[string]any) *map[string]any
}

// Extract parses commit fields to extract datas
func Extract(extractors *[]Extracter, commitMaps *[]map[string]any) *types.Changelog {
	results := []map[string]any{}

	for _, commitMap := range *commitMaps {
		result := &commitMap

		for _, extractor := range *extractors {
			result = extractor.Extract(result)
		}

		results = append(results, *result)
	}

	changelog := types.Changelog{}
	changelog.Datas = results
	changelog.Metadatas = map[string]any{}

	return &changelog
}

// Create builds extracters from a config
func Create(features Features, extractors Config) *[]Extracter {
	results := []Extracter{}

	if !features.ENABLED {
		return &results
	}

	for _, extractor := range extractors {
		results = append(results, regex{
			extractor.ORIGKEY,
			extractor.DESTKEY,
			extractor.REG,
		})
	}

	return &results
}
