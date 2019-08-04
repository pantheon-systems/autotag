APP := autotag

include scripts/make/common.mk
include scripts/make/common-go.mk

build::
	go build -o $(APP)/$(APP)  $(APP)/*.go

snapshot:
	@goreleaser --rm-dist --snapshot --debug
