[![Circle CI](https://circleci.com/gh/pantheon-systems/autotag.svg?style=shield&circle-token=ef9a68c180d0d470c594d39caf9e2a86fc529935)](https://circleci.com/gh/pantheon-systems/autotag)

# AutoTag

Automatically add version tags to a git repo based on commit messages.

Installing
------
Make sure you install 0.23.X of libgit on your box before trying to build this. The go2git lib is a wrapper around libgit.

Using the go tool you can get the cli with
```
  go get github.com/pantheon-systems/autotag/autotag

```

If you have `$GOPATH/bin` in your `$PATH` variable then you can run `autotag`  from the root of a git repo.


Usage
------
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
