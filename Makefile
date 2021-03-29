APP := autotag

include scripts/make/common.mk
include scripts/make/common-go.mk

build::
	CGO_ENABLED=0 go build -o $(APP)/$(APP)  $(APP)/*.go

snapshot:
	@goreleaser --rm-dist --snapshot --debug

.PHONY:: test-ci
test-ci:: test-go test-coveralls
