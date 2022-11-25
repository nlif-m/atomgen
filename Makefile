.POSIX:
.SUFFIX:
.PHONY: all vet fmt

GOFLAGS=-tags netgo -ldflags '-extldflags "-static"'


all: atomgen

atomgen: vet
	go build ${GOFLAGS}
vet: fmt
	go vet
fmt:
	go fmt
