FROM fnproject/go:dev as build-stage

ADD . /go/src/func
WORKDIR /go/src/func
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o  /hotwrap

FROM scratch

COPY --from=build-stage /hotwrap /hotwrap
