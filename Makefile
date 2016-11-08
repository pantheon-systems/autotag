appname=autotag
ARCH=$(shell uname -m)

include scripts/make/common.mk
include scripts/make/common-go.mk


# need to be able to commit with git to run tests on cirlce
deps-circle::
	git config --global user.email circleci
	git config --global user.name circleci
ifeq (, $(shell which gh-release))
	curl -L https://www.dropbox.com/s/4k3eq7xpehwwqr5/gh-release?dl=0 -o ~/bin/gh-release
	chmod 755 ~/bin/gh-release
endif

build::
	go build -o autotag/autotag  autotag/*.go

release:
	./autotag/autotag -n > VERSION
	mkdir release
	cp autotag/autotag release/autotag.linux.$(ARCH)
	gh-release create pantheon-systems/autotag $(shell ./autotag/autotag -n) $(shell git rev-parse --abbrev-ref HEAD)
