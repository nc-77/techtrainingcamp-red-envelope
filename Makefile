BUILD_DIR   := build
BUILD_FLAGS := -v

CGO_ENABLED := 0
GO111MODULE := on

LDFLAGS += -w -s -buildid=

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

GO_BUILD = GO111MODULE=$(GO111MODULE) CGO_ENABLED=$(CGO_ENABLED) \
	go build $(BUILD_FLAGS) -ldflags '$(LDFLAGS)' -trimpath

.PHONY: app clean

all: app

app:
	$(GO_BUILD) -o $(BUILD_DIR)/$@-$(GOOS)-$(GOARCH) main.go

clean:
	rm -rf $(BUILD_DIR)