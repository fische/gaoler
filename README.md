# Gaoler

A Go package manager

## TODO

**Clean this TODO list**

### Assumptions

* gaoler is executed in the directory where the `main` package is located

### Commands

* list:
  - list all imports (even subimports)
* save:
  - if it's not a repo, print warning
  - if it does not have any remote, print warning
* install:
  - go get in vendor/
  - checkout
  - remove .git, vendor and testdata of each repo
* update:
  - update on GOPATH
  - if not pinned, update according to changes in current GOPATH (check if commit is after)
  - if pinned, check if branch exist, store current state(branch/commit/whatever) for rolling back, checkout to branch, pull, rollback
* test:

### Config

* store only :
  - package name
  - commit
  - pinned branch

### Options

* dev dependencies
* scripts

### VCS

* Use package `git2go` with `libgit2` for the implementation of Git.

### Ideas
* use git sparse-checkout to save
