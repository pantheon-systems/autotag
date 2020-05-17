AutoTag
=======

[![Circle CI](https://circleci.com/gh/pantheon-systems/autotag.svg?style=shield&circle-token=ef9a68c180d0d470c594d39caf9e2a86fc529935)](https://circleci.com/gh/pantheon-systems/autotag)
[![Coverage Status](https://coveralls.io/repos/github/pantheon-systems/autotag/badge.svg?branch=master)](https://coveralls.io/github/pantheon-systems/autotag?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/pantheon-systems/autotag)](https://goreportcard.com/report/github.com/pantheon-systems/autotag)

Automatically increment version tags to a git repo based on commit messages.

Dependencies
------------

* [Git 2.x](https://git-scm.com/downloads) available in PATH

Version v1.0.0+ depends on the Git CLI, install Git with your distribution's package management
system.

Versions prior to v1.0.0 use cgo libgit or native golang Git, the binary will work standalone.

Installing
----------

### Pre-built binaries

| OS    | Arch  | binary              |
| ----- | ----- | ------------------- |
| macOS | amd64 | [autotag][releases] |
| Linux | amd64 | [autotag][releases] |

### Docker images

| Arch  | Images                                                           |
| ----- | ---------------------------------------------------------------- |
| amd64 | `quay.io/pantheon-public/autotag:latest`, `vX.Y.Z`, `vX.Y`, `vX` |

[releases]: https://github.com/pantheon-systems/autotag/releases/latest

### One-liner

Install Linux binary at `./autotag`. For example in a `Dockerfile` or a CI/CD pipeline:

```bash
curl -s https://api.github.com/repos/pantheon-systems/autotag/releases/latest | \
  grep browser_download | \
  grep Linux | \
  cut -d '"' -f 4 | \
  xargs curl -o ./autotag -L \
  && chmod 755 ./autotag
```

Usage
-----

The `autotag` utility will use the current state of the git repository to determine what the next
tag should be and then creates the tag by executing `git tag`. The `-n` flag will print the next tag but not apply it.

`autotag` scans the `master` branch for commits by default. Use `-b/--branch` to scan a different
branch. The utility first looks to find the most-recent reachable tag that matches a supported
versioning scheme. If no tags can be found the utility bails-out, so you do need to create a
`v0.0.0` tag before using `autotag`.

Once the last reachable tag has been found, the `autotag` utility inspects each commit between the
tag and `HEAD` of the branch to determine how to increment the version.

Commit messages are parsed for keywords via schemes. Schemes influence the tag selection according
to a set of rules.

Schemes are specified using the `-s/--scheme` flag:

### Scheme: Autotag (default)

The autotag scheme implements SemVer style versioning `vMajor.Minor.Patch` (e.g., `v1.2.3`).

Before using autotag for the first time create an initial SemVer tag,
eg: `git tag v0.0.0 -m'initial tag'`

The next version tag is calculated based on the contents of commit message according to these
rules:

- Bump the **major** version by including `[major]` or `#major` in a commit message, eg:

```
[major] breaking change
```

- Bump the **minor** version by including `[minor]` or `#minor` in a commit message, eg:

```
[minor] new feature added
```

- Bump the **patch** version by including `[patch]` or `#patch` in a commit message, eg:

```
[patch] bug fixed
```

If no keywords are specified a **Patch** bump is applied.

### Scheme: Conventional Commits

Specify the [Conventional Commits](TODO) v1.0.0 scheme by passing `--scheme=conventional` to `autotag`.

Conventional Commits implements SemVer style versioning `vMajor.Minor.Patch` similar to the
autotag scheme, but with a different commit message format.

Examples of Conventional Commits:

- A commit message footer containing `BREAKING CHANGE:` will bump the **major** version:

```
feat: allow provided config object to extend other configs

BREAKING CHANGE: `extends` key in config file is now used for extending other config files
```

- A commit message header containing a *type* of `feat` will bump the **minor** version:

```
feat(lang): add polish language
```

- A commit message header containg a `!` after the *type* is considered a breaking change and will
  bump the **major** version:

```
refactor!: drop support for Node 6
```

If no keywords are specified a **Patch** bump is applied.

### Pre-Release Tags

`autotag` supports appending additional test to the calculated next version string:

* Use `-p/--pre-release-name=` to append a pre-release **name** to the version. Allowed names are:
`alpha|beta|pre|rc|dev`. Example, `v1.2.3-dev`

* Use `-T/--pre-release-timestmap=` to append **timestamp** to the version. Allowed timetstamp
  formats are `datetime` (YYYYMMDDHHMMSS) or `epoch` (UNIX epoch timestamp in seconds).

Examples
--------

```console
$ autotag
3.2.1
```

```console
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

For additional help information use the `-h/--help` flag:

```console
autotag -h
```

### Goreleaser

`autotag` works well with [goreleaser](https://goreleaser.com/) for automating the process of
creating new versions and releases from CI.

An example of a [Circle-CI](https://circleci.com/) job utilizing both `autotag` and `goreleaser`:

```yaml
jobs:
  release:
    steps:
      - run:
          name: install autotag binary
          command: |
            curl -s https://api.github.com/repos/pantheon-systems/autotag/releases/latest | \
              grep browser_download | \
              grep -i linux | \
              cut -d '"' -f 4 | \
              xargs curl -o ~/autotag -L \
              && chmod 755 ~/autotag
      - run:
          name: increment version
          command: |
           ./autotag
      - run:
          name: build and push releases
          command: |
            curl -sL https://git.io/goreleaser | bash -s -- --parallelism=2 --rm-dist

workflows:
  version: 2
  build-test-release:
    jobs:
      - release
          requires:
            - build
          filters:
            branches:
              only:
                - master
```

Build from Source
-----------------

Assuming you have Go 1.5+ installed you can checkout and run make deps build to compile the binary
at `./autotag/autotag`.

```console
git clone git@github.com:pantheon-systems/autotag.git

cd autotag

make build
```

Release information
-------------------

Autotag itself uses `autotag` to increment releases. The default [autotag](#scheme-autotag-default) scheme is used for version selection.