[![Circle CI](https://circleci.com/gh/pantheon-systems/autotag.svg?style=shield&circle-token=ef9a68c180d0d470c594d39caf9e2a86fc529935)](https://circleci.com/gh/pantheon-systems/autotag)
[![Coverage Status](https://coveralls.io/repos/github/pantheon-systems/autotag/badge.svg?branch=master)](https://coveralls.io/github/pantheon-systems/autotag?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/pantheon-systems/autotag)](https://goreportcard.com/report/github.com/pantheon-systems/autotag)

AutoTag
-------

Automatically add version tags to a git repo based on commit messages.

Installing
==========

On Linux the easy way to get going is to use the pre-built binary release from github releases. 


Usage
=====

The default behavior with no arguments will tag a new version on current repo and emit the version tagged
```
$ autotag
v3.2.1
```

you can get more help using -h flag
```
$ autotag -h
Usage:
  autotag [OPTIONS]

Application Options:
  -n          Just output the next version, don't autotag
  -v          Enable verbose logging
  -r, --repo= Path to the repo (./)

Help Options:
  -h, --help  Show this help message
```

Build from Source
=================
Assuming you have Go 1.5+ installed you can checkout and run make deps build to compile the binary. It will be built as ./autotag/autotag


```
git clone git@github.com:pantheon-systems/autotag.git 

cd autotag

make deps build
```
