.POSIX:
.SUFFIX:
.PHONY: all vet fmt windows_amd64 linux_arm linux_amd64 crossplatform atomgen
PROGRAM_NAME=atomgen
GOOS=linux
GOARCH=amd64
OUTPUT_FILE="${PROGRAM_NAME}_${GOOS}_${GOARCH}"
GOFLAGS=-tags netgo -ldflags '-extldflags "-static"'


all: atomgen

atomgen: vet test
	GOOS=${GOOS} GOARCH=${GOARCH} go build ${GOFLAGS} -o "$@"

vet: fmt
	go vet ./...

fmt:
	go fmt ./...
	goimports -l -w .

test:
	go test ./...

linux_amd64: vet test
	GOOS=linux GOARCH=amd64 go build ${GOFLAGS} -o "${PROGRAM_NAME}_linux_amd64"

# windows_amd64: vet test
# 	GOOS=windows GOARCH=amd64 go build ${GOFLAGS} -o "${PROGRAM_NAME}_windows_amd64.exe"

# linux_arm: vet test
# 	GOOS=linux GOARCH=arm go build ${GOFLAGS} -o "${PROGRAM_NAME}_linux_arm"

# crossplatform: linux_amd64 linux_arm windows_amd64 
