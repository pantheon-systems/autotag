[![Circle CI](https://circleci.com/gh/pantheon-systems/autotag.svg?style=shield&circle-token=ef9a68c180d0d470c594d39caf9e2a86fc529935)](https://circleci.com/gh/pantheon-systems/autotag)

# AutoTag

Automatically add version tags to a git repo based on commit messages.

## Installing
On Linux the easy way to get going is to use the prebuilt binary relese from github releases. This binary has a static version of libgit2 embedded in it, and should be fairly portable to most linux distros.

If you are not on linux/x86_64 you will have to build it from source via the instructions below.

## Usage
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

## Build from Source
If you want to build from source you will need to have a version of libgit2  >= 0.23.0 installed on your system before hand. The go2git lib is a wrapper around libgit.

After having installed libgit2 you can install the cli using the go tool:
```
  go get github.com/pantheon-systems/autotag/autotag

```

If you have `$GOPATH/bin` in your `$PATH` variable then you can run `autotag`  from the root of a git repo.

