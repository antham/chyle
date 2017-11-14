package prompt

import (
	"github.com/antham/strumt"

	"github.com/antham/chyle/prompt/internal/builder"
)

func newExtractors(store *builder.Store) []strumt.Prompter {
	return builder.NewGroupEnvPromptWithCounter(extractor, store)
}

var extractor = []builder.EnvConfig{
	{"extractorOrigKey", "extractorDestKey", "CHYLE_EXTRACTORS_*_ORIGKEY", "Enter a commit field from which we want to extract datas (id, authorName, authorEmail, authorDate, committerName, committerEmail, committerMessage, type)"},
	{"extractorDestKey", "extractorReg", "CHYLE_EXTRACTORS_*_DESTKEY", "Enter a name for the key which will receive the extracted value"},
	{"extractorReg", "mainMenu", "CHYLE_EXTRACTORS_*_REG", "Enter a regexp used to extract a data"},
}
