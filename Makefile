.POSIX:
.SUFFIX:
.PHONY: all vet fmt

GOFLAGS=-tags netgo -ldflags '-extldflags "-static"'


all: atomgen

atomgen: vet test
	go build ${GOFLAGS}
vet: fmt
	go vet
fmt:
	go fmt

test:
	go test
