[![Circle CI](https://circleci.com/gh/pantheon-systems/autotag.svg?style=shield&circle-token=ef9a68c180d0d470c594d39caf9e2a86fc529935)](https://circleci.com/gh/pantheon-systems/autotag)
[![Coverage Status](https://coveralls.io/repos/github/pantheon-systems/autotag/badge.svg?branch=master)](https://coveralls.io/github/pantheon-systems/autotag?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/pantheon-systems/autotag)](https://goreportcard.com/report/github.com/pantheon-systems/autotag)

AutoTag
-------

Automatically add version tags to a git repo based on commit messages.

Dependencies
==========
* [Git 2.x](https://git-scm.com/downloads) available in PATH

Installing
==========

On Linux the easy way to get going is to use the pre-built binary release from [GitHub Releases](https://github.com/pantheon-systems/autotag/releases). 

If using a recent version that depends on the Git CLI, install Git with your distribution's package management system. 

If using an older release with cgo libgit or native golang Git, the binary will work standalone.

Calculating Tags
================

The `autotag` utility will use the current state of the git repository to determine what the next tag should be (when following SemVer 2.0).
Tags created by `autotag` have the following format: `vMajor.Minor.Patch` (e.g., `v1.2.3`).

By default, `autotag` only scans the `master` branch for changes. The utility first looks to find the most-recent reachable tag, only
looking for tags that appear to be version strings. If no tags can be found the utility bails-out, so you do need to create a `v0.0.0` tag
before using `autotag`.

Once the last reachable tag has been found, the `autotag` utility inspects each commit between the tag and `HEAD` of the branch to determine
how to increment the version. By default a single `Patch` increase is made (i.e., `v1.2.3` => `v1.2.4`). However, information can be included
in a commit to tell `autotag` to increment the `Major` and `Minor` versions.

### Incrementing Major and Minor versions

When the `autotag` utility inspects the commits between the latest tag and `HEAD`, it looks for certain strings to tell it to increment
something other than the `Patch` version. This is a simple regular expression match against your commit message.

To increase your `Major` version, you can include either `[major]` or `#major` in your commit message. That means you can have the subject
of your commit be:

```
[major] version bump in preparation for release
```

Likewise, you can include them anywhere in your commit message:

```
Fix the thing with the stuff

WISOTT
#major
```

This would result in `v1.2.3` becoming `v2.0.0`. Telling `autotag` to increase the `Minor` version is the same as with the `Major`, except
use `[minor]` or `#minor` instead. A `Minor` version bump would result in a change from `v1.2.3` to `v1.3.0`.

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
