package chyle

// decorater extends data from commit hashmap with data picked from third part service
type decorater interface {
	decorate(*map[string]interface{}) (*map[string]interface{}, error)
}

// decorate process all defined decorator and apply them
func decorate(decorators *map[string][]decorater, changelog *Changelog) (*Changelog, error) {
	var err error

	datas := []map[string]interface{}{}

	for _, d := range changelog.Datas {
		result := &d

		for _, decorator := range (*decorators)["datas"] {
			result, err = decorator.decorate(&d)

			if err != nil {
				return nil, err
			}
		}

		datas = append(datas, *result)
	}

	changelog.Datas = datas

	metadatas := changelog.Metadatas

	for _, decorator := range (*decorators)["metadatas"] {
		m, err := decorator.decorate(&metadatas)

		if err != nil {
			return nil, err
		}

		metadatas = *m
	}

	changelog.Metadatas = metadatas

	return changelog, nil
}

// createDecorators build decorators from a config
func createDecorators() *map[string][]decorater {
	results := map[string][]decorater{"metadatas": {}, "datas": {}}

	if chyleConfig.FEATURES.HASJIRADECORATOR {
		results["datas"] = append(results["datas"], buildJiraDecorator())
	}

	if chyleConfig.FEATURES.HASENVDECORATOR {
		results["metadatas"] = append(results["metadatas"], buildEnvDecorators()...)
	}

	return &results
}
