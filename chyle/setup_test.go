package chyle

import (
	"os"
	"os/exec"
	"testing"

	"github.com/Sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
)

var repo *git.Repository

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	for _, filename := range []string{"../features/init.sh", "../features/merge-commits.sh"} {
		err := exec.Command(filename).Run()

		if err != nil {
			logrus.Fatal(err)
		}
	}

	path, err := os.Getwd()

	if err != nil {
		logrus.Fatal(err)
	}

	repo, err = git.NewFilesystemRepository(path + "/test/.git/")

	if err != nil {
		logrus.Fatal(err)
	}
}
