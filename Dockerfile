FROM golang:alpine as builder
RUN apk add --no-cache git
COPY . $GOPATH/src/paam/
WORKDIR $GOPATH/src/paam/
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/paam

FROM scratch
COPY --from=builder /go/bin/paam /go/bin/paam
EXPOSE 8113
ENTRYPOINT ["/go/bin/paam"]
