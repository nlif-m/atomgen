from docker.io/library/golang:1.19-alpine as builder
RUN apk add --no-cache build-base && \
    go install golang.org/x/tools/cmd/goimports@v0.19.0

COPY . /app
WORKDIR /app
RUN make

from docker.io/library/alpine:3.18 
WORKDIR /app

VOLUME atomgen-config
VOLUME src
EXPOSE 3000

RUN apk add --no-cache python3 py3-pip ffmpeg && \
    pip3 --no-cache-dir install yt-dlp

COPY --from=builder /app/atomgen /bin/atomgen

ENTRYPOINT ["/bin/atomgen", "-config", "/etc/atomgen/config.json"]



