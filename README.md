# Gaoler

A Go package manager

## V0.1

* [ ] **clean code (split it into more functions and modules)**
* [ ] find `main` package automatically

### Commands

- [ ] vendor:
  - if it's not a repo, print warning
  - if it does not have any remote, print warning
  - by default, it reads `gaoler.json` to retrieve all dependencies
  - add option to regenerate list of dependencies and save it to `gaoler.json`
- [ ] update:
  - update on GOPATH
  - if not pinned, update according to changes in current GOPATH (check if commit is after)
  - if pinned, check if branch exist, store current state(branch/commit/whatever) for rolling back, checkout to branch, pull, rollback

### Config

* [ ] store only :
  - package name
  - remotes
  - commit
  - pinned branch

### VCS

* [x] Module :
  - git
* [ ] Use package `git2go` with `libgit2` for the implementation of Git.

### Options

* [ ] dev dependencies

## Suggestion box

### Commands

* restore:
  - restore GOPATH using list of dependencies from `gaoler.json`

### VCS

* Modules :
  - hg
  - bzr
  - svn
* use git sparse-checkout to retrieve only files that we need

### Config

* add scripts
