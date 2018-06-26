FROM golang:1.10 as builder

COPY vendor       /go/src/github.com/fnproject/hotwrap/vendor
COPY  hotwrap.go  /go/src/github.com/fnproject/hotwrap/hotwrap.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o  /hotwrap  /go/src/github.com/fnproject/hotwrap/hotwrap.go

FROM scratch

COPY --from=builder /hotwrap /hotwrap
