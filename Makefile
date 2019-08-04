APP := autotag

include scripts/make/common.mk
include scripts/make/common-go.mk

build::
	go build -o $(APP)/$(APP)  $(APP)/*.go

release: VERSION=$(shell $(APP)/$(APP) -n)
release:
	GOOS=darwin go build -o $(APP)/$(APP)-darwin autotag/*.go
	github-release release -u pantheon-systems -r $(APP) -t $(VERSION) --draft
	github-release upload -u pantheon-systems -r $(APP) -n Linux -f $(APP)/$(APP) -t $(VERSION)
	github-release upload -u pantheon-systems -r $(APP) -n OSX -f $(APP)/$(APP)-darwin -t $(VERSION)
