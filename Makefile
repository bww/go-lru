
export GOPATH := $(GOPATH):$(PWD)

SRC=src/lru/*.go

.PHONY: all deps test

all: test

deps:

test:
	go test lru -test.v

