package prompt

import (
	"github.com/antham/strumt"
)

func newExtractors(store *Store) []strumt.Prompter {
	return newGroupEnvPromptWithCounter(extractor, store)
}

var extractor = []envConfig{
	envConfig{"extractorOrigKey", "extractorDestKey", "CHYLE_EXTRACTORS_*_ORIGKEY", "Enter a commit field from which we want to extract datas (id, authorName, authorEmail, authorDate, committerName, committerEmail, committerMessage, type)"},
	envConfig{"extractorDestKey", "extractorReg", "CHYLE_EXTRACTORS_*_DESTKEY", "Enter a name for the key which will receive the extracted value"},
	envConfig{"extractorReg", "mainMenu", "CHYLE_EXTRACTORS_*_REG", "Enter a regexp used to extract a data"},
}
