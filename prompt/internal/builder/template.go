package builder

type template struct {
	iD        string
	PromptStr string
	onSuccess func(string) string
	onError   func(error) string
	parse     func(string) error
}

func (t *template) ID() string {
	return t.iD
}

func (t *template) PromptString() string {
	return t.PromptStr
}

func (t *template) Parse(value string) error {
	return t.parse(value)
}

func (t *template) NextOnSuccess(value string) string {
	return t.onSuccess(value)
}

func (t *template) NextOnError(err error) string {
	return t.onError(err)
}
