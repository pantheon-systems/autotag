all: test build cover
godir="${HOME}/go1.5/src/github.com/pantheon-systems"
appname=autotag
ARCH=$(shell uname -m)

# check for circle, and our repo in the gopath on circle. make a link if it's not there
deps:
	go get github.com/axw/gocov/gocov
	go get github.com/mattn/goveralls
	go get github.com/golang/lint/golint
	go get github.com/hashicorp/go-version
	go get github.com/jessevdk/go-flags
	go get -d github.com/libgit2/git2go
ifdef CIRCLECI
	scripts/static-git2go.sh
	test -s ${godir} || mkdir -p ${godir}
	test -s ${godir}/${appname} || ln -s ${HOME}/${appname} ${godir}/${appname}
endif

build: deps
	go build -o autotag/autotag  autotag/*.go

test: deps
	test -z "$(gofmt -s -l . | tee /dev/stderr)"
	go vet .
	test -z "$(golint . | tee /dev/stderr)"
	go test -v .

cover:
	go test -cover

cov:
	gocov test ./... | gocov-html > /tmp/coverage.html
	open /tmp/coverage.html

release:
	./autotag/autotag -n > VERSION
	mkdir release
	tar -zcf release/autotag-linux.$(ARCH).tgz autotag/autotag
	gh-release create pantheon-systems/autotag $(shell ./autotag/autotag -n) $(shell git rev-parse --abbrev-ref HEAD)

.PHONY: all cov test
