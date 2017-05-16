package decorators

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/antham/chyle/chyle/convh"
)

type shellConfig map[string]struct {
	COMMAND string
	ORIGKEY string
	DESTKEY string
}

// shell pipes a shell command on field content and dump the
// result into a new field
type shell struct {
	COMMAND string
	ORIGKEY string
	DESTKEY string
}

// Decorate executes shell command on field content
func (s shell) Decorate(commitMap *map[string]interface{}) (*map[string]interface{}, error) {
	var tmp interface{}
	var value string
	var result []byte
	var ok bool
	var err error

	(*commitMap)[s.DESTKEY] = nil

	if tmp, ok = (*commitMap)[s.ORIGKEY]; !ok {
		return commitMap, nil
	}

	if value, err = convh.ConvertToString(tmp); err != nil {
		return commitMap, nil
	}

	command := fmt.Sprintf(`echo "%s"|%s`, strings.Replace(value, `"`, `\"`, -1), s.COMMAND)

	if result, err = exec.Command("sh", "-c", command).Output(); err != nil { // #nosec
		return commitMap, fmt.Errorf("%s : command failed", command)
	}

	(*commitMap)[s.DESTKEY] = string(result[:len(result)-1])

	return commitMap, nil
}

// buildShell create a new shell decorator
func buildShell(configs shellConfig) []Decorater {
	results := []Decorater{}

	for _, config := range configs {
		results = append(results, shell{config.COMMAND, config.ORIGKEY, config.DESTKEY})
	}

	return results
}
