# Gaoler

A Go dependency manager

## Usage

```
Usage: goaler [-v] [--config=<config-file>] [--main=<main-package>] COMMAND [arg...]

A Go dependency manager

Options:
  -m, --main="."                 Path to the main package
  -c, --config="./gaoler.json"   Path to the configuration file
  -v, --verbose=false            Enable verbose mode

Commands:
  list         List dependencies of your project
  update       Update dependencies of your project
  vendor       Vendor dependencies of your project

Run 'goaler COMMAND --help' for more information on a command.
```

## Global Options and Arguments

* `--config` defines the location of the config file
* `--main` defines the location of the root directory of the project including the package main

## Commands

* Find `main` package automatically
* List:
  - Lists all dependencies of the project, even local ones
  - Include test dependencies with option `-t`
* Vendor:
  - Vendor imported packages
  - Include test dependencies with option `-t`
  - Save depedencies with option `-s` to the config file
  - Load by default the config file to only vendor new packages (Use option `-f` to not use it)
  - Clean unnecessary directories (.git and unused directories)
* Update:
  - Update based branch updates

## VCS

* Module :
  - git

## Config

* You can choose the format between YAML and JSON by changing the extension of the config file in the options
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

## Next release (0.4)

* Unit tests
* E2E tests
* Write a clean documentation

### Packages

* Logger
* Error

## Suggestion box

* Clean README and define conventions

### Config

* Split test dependencies from the normal ones

### Commands

* Clean:
  - clean vendor directory from any unused dependency
* Restore:
  - Restore vendor to statement as indicated in `gaoler.json`
* Update:
  - update based on GOPATH (used for packages without any remote repository or if you only want to update using GOPATH instead of the remote)
  - update using tags and similar syntax as npm

### VCS

* Modules :
  - hg
  - bzr
  - svn
* Use git sparse-checkout to retrieve only files that we need
* Support whatever remote name

* **Pending:** Use package `git2go` with `libgit2` for the implementation of Git.

### Options

* Scripts
