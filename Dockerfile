FROM golang:1.6.2-alpine

RUN apk add --update git && apk add --update make && rm -rf /var/cache/apk/*

ADD . /go/src/github.com/r3labs/instance-adapter
WORKDIR /go/src/github.com/r3labs/instance-adapter

RUN make deps && go install

ENTRYPOINT /go/bin/instance-adapter
