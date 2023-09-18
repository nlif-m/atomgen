from docker.io/library/golang:1.19-alpine as builder
RUN apk update && \
    apk upgrade &&  \
    apk add build-base


COPY . /app
WORKDIR /app
RUN make

from docker.io/library/alpine:3.18 as runner
WORKDIR /app

VOLUME atomgen-config
VOLUME src

RUN apk update && \
    apk upgrade && \
    apk add python3 py3-pip && \
    pip3 --no-cache-dir install yt-dlp

COPY --from=builder /app/atomgen /bin/atomgen

ENTRYPOINT ["atomgen"]



