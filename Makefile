APP=autotag

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
	go build -o $(APP)/$(APP)  $(APP)/*.go

release: VERSION=$(shell $(APP)/$(APP) -n)
release:
	GOOS=darwin go build -o $(APP)/$(APP)-darwin autotag/*.go
	github-release release -u pantheon-systems -r $(APP) -t $(VERSION) --draft
	github-release upload -u pantheon-systems -r $(APP) -n Linux -f $(APP)/$(APP) -t $(VERSION)
	github-release upload -u pantheon-systems -r $(APP) -n OSX -f $(APP)/$(APP)-darwin -t $(VERSION)
