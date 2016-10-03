package vcs_test

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

type RepositoryTest struct {
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

func (r RepositoryTest) GetRevision() (string, error) {
	return r.revision()
}

func (r RepositoryTest) GetRemote() (string, error) {
	return r.remote()
}

func (r RepositoryTest) GetVCSName() string {
	return r.name
}

func (r RepositoryTest) GetPath() (string, error) {
	return r.path()
}

func (r RepositoryTest) GetBranch() (string, error) {
	return r.branch()
}

func (r RepositoryTest) AddRemote(remote string) error {
	return r.addRemote(remote)
}

func (r RepositoryTest) Fetch() error {
	return r.fetch()
}

func (r RepositoryTest) Checkout(revision string) error {
	return r.checkout(revision)
}

func (r RepositoryTest) CheckoutBranch(branch string) error {
	return r.checkoutBranch(branch)
}

func TestVcs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VCS Suite")
}
