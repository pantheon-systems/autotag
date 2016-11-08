appname=autotag

include scripts/make/common.mk
include scripts/make/common-go.mk

# need to be able to commit with git to run tests on cirlce
deps-circle::
	git config --global user.email circleci
	git config --global user.name circleci
ifeq (, $(shell which gihub-release))
	go get github.com/aktau/github-release
endif

build::
	go build -o autotag/autotag  autotag/*.go

release: VERSION=$(shell $(AUTOTAG) -n)
release:
	mkdir release
	cp autotag/autotag autotag
	GOOS=darwin go build -o $(APP)-darwin autotag/*.go
	github-release create pantheon-systems/autotag $(shell ./autotag/autotag -n) $(shell git rev-parse --abbrev-ref HEAD)
	github-release release -u pantheon-systems -r $(APP) -t $(VERSION) --draft
	github-release upload -u pantheon-systems -r $(APP) -n Linux -f $(APP) -t $(VERSION)
	github-release upload -u pantheon-systems -r $(APP) -n OSX -f $(APP)-darwin -t $(VERSION)
