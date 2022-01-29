FROM golang:1.15-buster

WORKDIR /go/src/app
RUN go get github.com/cespare/reflex

CMD ["make","watcher"]
