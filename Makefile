.POSIX:
.SUFFIX:
.PHONY: all vet fmt

all: atomgen

atomgen: vet
	go build

vet: fmt
	go vet
fmt:
	go fmt
