BASIC=Makefile
GO=go
GIT=git

$(eval VERSION:=$(shell git rev-parse HEAD))
LD_FLAGS=-ldflags "-X main.Version=${VERSION}"

all: .PHONY ingest

ingest: .PHONY
	$(GO) build $(LD_FLAGS)

fmt: .PHONY
	$(GO) fmt

test: .PHONY

.PHONY:
