# Gaoler

A Go package manager

## V0.1

### Commands

- [x] vendor:
  - clean unnecessary directories

### VCS

* [x] Module :
  - git
* [ ] Use package `git2go` with `libgit2` for the implementation of Git.

## Suggestion box

### Logger - Error Packages

### Config

* Proposition :

```json
{
  "dependencies": {
    "github.com/fische/gaoler": {
      "remote": "https://github.com/fische/gaoler",
      "vcs": "git",
      "revision": "revision",
      "branch": "master",
      "packages": [
        "github.com/fische/gaoler",
        "github.com/fische/gaoler/cmd",
        "github.com/fische/gaoler/errors",
        "github.com/fische/gaoler/project",
        "github.com/fische/gaoler/vcs",
        "github.com/fische/gaoler/vcs/modules",
        "github.com/fische/gaoler/vcs/modules/git",
      ]
    }
  }
}
```

### Commands

* Vendor:
  - by default, it reads `gaoler.json` to retrieve all dependencies
  - add option to regenerate list of dependencies and save it to `gaoler.json`
* Restore:
  - Restore GOPATH using list of dependencies from `gaoler.json`
* Update:
  - update on GOPATH
  - if not pinned, update according to changes in current GOPATH (check if commit is after)
  - if pinned, check if branch exist, store current state(branch/commit/whatever) for rolling back, checkout to branch, pull, rollback


### VCS

* Modules :
  - hg
  - bzr
  - svn
* Use git sparse-checkout to retrieve only files that we need
* Support whatever remote name

### Options

* Dev dependencies
* Scripts
* Find `main` package automatically
