package git_test

import (
	"log"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func resetDirectory(dir string) {
	if err := os.RemoveAll(dir); err != nil && !os.IsNotExist(err) {
		log.Fatalf("%v", err)
	} else if err = os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("%v", err)
	}
}

func removeDirectory(dir string) {
	if err := os.RemoveAll(dir); err != nil && !os.IsNotExist(err) {
		log.Fatalf("%v", err)
	}
}

func TestGit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Git Suite")
}
