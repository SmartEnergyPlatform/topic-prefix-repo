FROM golang:1.11


COPY . /go/src/topic-prefix-repo
WORKDIR /go/src/topic-prefix-repo

ENV GO111MODULE=on

RUN go build

EXPOSE 8080

CMD ./topic-prefix-repo