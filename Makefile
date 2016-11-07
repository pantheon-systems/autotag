appname=autotag
ARCH=$(shell uname -m)

include scripts/make/common.mk
include scripts/make/common-go.mk


# need to be able to commit with git to run tests on cirlce
deps-circle::
	git config --global user.email circleci
	git config --global user.name circleci

build::
	go build -o autotag/autotag  autotag/*.go

release:
	./autotag/autotag -n > VERSION
	mkdir release
	cp autotag/autotag release/autotag.linux.$(ARCH)
	gh-release create pantheon-systems/autotag $(shell ./autotag/autotag -n) $(shell git rev-parse --abbrev-ref HEAD)
