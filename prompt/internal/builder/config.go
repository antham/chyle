package builder

// EnvConfig is common config for all environments variables prompts builder
type EnvConfig struct {
	ID           string
	NextID       string
	Env          string
	PromptString string
	Validator    func(value string) error
}