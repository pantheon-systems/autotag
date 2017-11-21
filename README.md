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
something other than the `Patch` version. This is a simple regular expression match against your commit message subject. If you don't
place the marker in your commit subject line, `autotag` will not observe it and won't correctly bump the version.

To increase your `Major` version, you can include either `[major]` or `#major` in your commit message. That means you can have the subject
of your commit be:

```
[major] version bump in preparation for release
```

Or if you prefer hashtags:

```
#major version bump in preparation for release
```

This would result in `v1.2.3` becoming `v2.0.0`. Telling `autotag` to increase the `Minor` version is the same as with the `Major`, except
use `[minor]` or `#minor` instead. A `Minor` version bump would result in a change from `v1.2.3` to `v1.3.0`.

### Pre-Release Tags

The `autotag` package supports providing a `PreReleaseName` and a `PreReleaseTimestampLayout`, which gives you the ability to automatically
create tags like `v1.2.3-pre.20170706070042`. You can omit the name, or the timestamp, to include as much information as you'd like. The
timestamp layout value uses the standard layout from the `time` package.

This works by finding the last tag without any pre-release information, say `v1.2.3`. It then bumps the version based on the rules above,
and appends our pre-release information on to the end starting with a hyphen. For example:

* `-<PreReleaseName>`
* `-<PreReleaseTimestamp>`
* `-<PreReleaseName>.<PreReleaseTimestamp>`

The `autotag` binary provides controlled access to this functionality by allowing you to choose one of four pre-release names, as well
as by either using a UNIX `epoch` timestamp or a `datetime` timestamp in the form of `YYYYMMDDHHMMSS`. The `pre-release-name` is
implemented in the `-p` flag, while the timestamp layout is implemented in the `-T` flag. See the [Usage](#Usage) section for more
information.

Usage
=====

The default behavior with no arguments will tag a new version on current repo and emit the version tagged:

```
$ autotag
3.2.1
```

`autotag` also supports pre-release tags with the `-p` and `-T` flags, and here are some example:

```
$ autotag -p pre
3.2.1-pre

$ autotag -T epoch
3.2.1-1499320004

$ autotag -T datetime
3.2.1-20170706054703

$ autotag -p pre -T epoch
3.2.1-pre.1499319951

$ autotag -p rc -T datetime
3.2.1-rc.20170706054528
```


You can get more help using the `-h/--help` flag:

```
$ autotag -h
Usage:
  autotag [OPTIONS]

Application Options:
  -n                           Just output the next version, don't autotag
  -v                           Enable verbose logging
  -b, --branch=                Git branch to scan (default: master)
  -r, --repo=                  Path to the repo (default: ./)
  -p, --pre-release-name=      create a pre-release tag with this name (can be: alpha|beta|rc|pre|hotfix)
  -T, --pre-release-timestamp= create a pre-release tag and append a timestamp (can be: datetime|epoch)

Help Options:
  -h, --help                   Show this help message
```

Build from Source
=================
Assuming you have Go 1.5+ installed you can checkout and run make deps build to compile the binary. It will be built as ./autotag/autotag


```
git clone git@github.com:pantheon-systems/autotag.git 

cd autotag

make deps build
```
