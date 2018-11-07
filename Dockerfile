FROM fnproject/go:dev as build-stage

RUN apk update && apk add bash
ADD . /go/src/func
WORKDIR /go/src/func
RUN go test -v ./... && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o  /hotwrap

FROM fnproject/go

COPY --from=build-stage /hotwrap /hotwrap
