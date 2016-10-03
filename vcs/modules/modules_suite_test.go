package modules_test

import (
	"github.com/fische/gaoler/vcs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

type VCSTest struct {
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

func (v VCSTest) New(path string) (vcs.Repository, error) {
	return v.new(path)
}

func (v VCSTest) Open(path string) (vcs.Repository, error) {
	return v.open(path)
}

func (v VCSTest) GetName() string {
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

func TestModules(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Modules Suite")
}
