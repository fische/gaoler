# Gaoler

A Go package manager

## Usage

```
Usage: goaler [-v] [--config=<config-file>] [ROOT] COMMAND [arg...]

A Go package manager

Arguments:
  ROOT="/home/fische_m/Project/Go/src/github.com/fische/gaoler"   Root directory from a project

Options:
  -c, --config="gaoler.json"   Path to the configuration file
  -v, --verbose=false          Enable verbose mode

Commands:
  list         List dependencies of your project
  vendor       Vendor dependencies of your project

Run 'goaler COMMAND --help' for more information on a command.
```

## Global Options and Arguments

* `-c` defines the location of the config file
* `ROOT` defines the location of the root directory of the project including the package main

## Commands

* List:
  - Lists all dependencies of the project, even local ones
  - Include test dependencies with option `-t`
* Vendor:
  - Vendor imported packages
  - Include test dependencies with option `-t`
  - Save depedencies with option `-s` to the config file
  - Load by default the config file to only vendor new packages (Use option `-f` to not use it)
  - Clean unnecessary directories (.git and unused directories)

## VCS

* Module :
  - git

### Config

* Scheme :

```JSON
{
  "Name": "github.com/fische/gaoler",
  "Dependencies": {
    "github.com/Sirupsen/logrus": {
      "Remote": "https://github.com/Sirupsen/logrus",
      "VCS": "git",
      "Revision": "revision",
      "Branch": "master",
      "Packages": [
        "github.com/Sirupsen/logrus"
      ]
    }
  }
}
```

## Next release (0.3)

### Config

* Let the user choose between json and yaml to save the config

### Commands

* Find `main` package automatically
* Update:
  - if not pinned, update according to changes in current GOPATH (check if commit is after)
  - if pinned, check if branch exist, store current state(branch/commit/whatever) for rolling back, checkout to branch, pull, rollback

## Suggestion box

* Write a clean documentation
* Write tests for the whole projects (unit tests as well as E2E)

### Packages

* Logger
* Error

### Commands

* Clean:
  - clean vendor directory from any unused dependency
* Restore:
  - Restore GOPATH using list of dependencies from `gaoler.json`

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
