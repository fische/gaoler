# Gaoler

A Go package manager

## Implemented

### Commands

* List:
  - Lists all dependencies of the project, even local ones
* Vendor:
  - Vendor all dependencies of the project (not local ones)
  - Clean unnecessary directories (.git and unused directories)

### VCS

* Module :
  - git

## Next release (0.2)

### Config

* Unmarshal `Project` from `gaoler.json` (default)
* Marshal `Project` to `gaoler.json` (default)
* Add option to specify config path
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
  - Add possibility to only vendor new imported packages

## Suggestion box

### Packages

* Logger
* Error

### Commands

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

* **Pending:** Use package `git2go` with `libgit2` for the implementation of Git.

### Options

* Dev dependencies
* Scripts
* Find `main` package automatically
