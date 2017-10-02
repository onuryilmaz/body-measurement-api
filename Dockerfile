#FROM golang:1.9 as tester
#ADD . /go/src/github.com/onuryilmaz/body-measurement-api/
#WORKDIR /go/src/github.com/onuryilmaz/body-measurement-api
#RUN GOOS=linux go test -v ./...

FROM golang:1.9 as builder
ADD . /go/src/github.com/onuryilmaz/body-measurement-api/
WORKDIR /go/src/github.com/onuryilmaz/body-measurement-api/cmd/data-api
RUN GOOS=linux go build -o data-api
WORKDIR /go/src/github.com/onuryilmaz/body-measurement-api/cmd/tracking-api
RUN GOOS=linux go build -o tracking-api

FROM ubuntu as data-api
COPY --from=builder /go/src/github.com/onuryilmaz/body-measurement-api/cmd/data-api/data-api /data-api
ENTRYPOINT ["./data-api"]

FROM ubuntu as tracking-api
COPY --from=builder /go/src/github.com/onuryilmaz/body-measurement-api/cmd/tracking-api/tracking-api /tracking-api
ENTRYPOINT ["./tracking-api"]