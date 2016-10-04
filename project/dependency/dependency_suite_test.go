package dependency_test

import (
	"log"
	"os"
	"time"

	"github.com/fische/gaoler/vcs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

type fileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
	sys     interface{}
}

type VCS struct {
	name string
	new  func(path string) (vcs.Repository, error)
	open func(path string) (vcs.Repository, error)
}

type Repository struct {
	name     string
	revision func() (string, error)
	remote   func() (string, error)
	path     func() (string, error)
	branch   func() (string, error)

	addRemote      func(remote string) error
	fetch          func() error
	checkout       func(revision string) error
	checkoutBranch func(branch string) error
}

func (v VCS) New(path string) (vcs.Repository, error) {
	return v.new(path)
}

func (v VCS) Open(path string) (vcs.Repository, error) {
	return v.open(path)
}

func (v VCS) GetName() string {
	return v.name
}

func (r Repository) GetRevision() (string, error) {
	return r.revision()
}

func (r Repository) GetRemote() (string, error) {
	return r.remote()
}

func (r Repository) GetVCSName() string {
	return r.name
}

func (r Repository) GetPath() (string, error) {
	return r.path()
}

func (r Repository) GetBranch() (string, error) {
	return r.branch()
}

func (r Repository) AddRemote(remote string) error {
	return r.addRemote(remote)
}

func (r Repository) Fetch() error {
	return r.fetch()
}

func (r Repository) Checkout(revision string) error {
	return r.checkout(revision)
}

func (r Repository) CheckoutBranch(branch string) error {
	return r.checkoutBranch(branch)
}

func getter(v string, err error) func() (string, error) {
	return func() (string, error) {
		return v, err
	}
}

func (f fileInfo) Name() string {
	return f.name
}

func (f fileInfo) Size() int64 {
	return f.size
}

func (f fileInfo) Mode() os.FileMode {
	return f.mode
}

func (f fileInfo) ModTime() time.Time {
	return f.modTime
}

func (f fileInfo) IsDir() bool {
	return f.isDir
}

func (f fileInfo) Sys() interface{} {
	return f.sys
}

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

func resetFile(file string) {
	if err := os.Remove(file); err != nil && !os.IsNotExist(err) {
		log.Fatalf("%v", err)
	} else if _, err = os.Create(file); err != nil {
		log.Fatalf("%v", err)
	}
}

func removeFile(file string) {
	if err := os.Remove(file); err != nil && !os.IsNotExist(err) {
		log.Fatalf("%v", err)
	}
}

func TestDependency(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dependency Suite")
}
