package chyle

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/antham/envh"
	"gopkg.in/src-d/go-git.v4"
)

var envs map[string]string
var repo *git.Repository

func TestMain(m *testing.M) {
	saveExistingEnvs()
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	for _, filename := range []string{"../features/init.sh", "../features/merge-commits.sh"} {
		err := exec.Command(filename).Run()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	path, err := os.Getwd()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	repo, err = git.NewFilesystemRepository(path + "/test/.git/")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func saveExistingEnvs() {
	var err error
	env := envh.NewEnv()

	envs, err = env.FindEntries(".*")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func restoreEnvs() {
	os.Clearenv()

	if len(envs) != 0 {
		for key, value := range envs {
			setenv(key, value)
		}
	}
}

func setenv(key string, value string) {
	err := os.Setenv(key, value)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
