FROM golang:1.8.1 as builder
ADD . /go/src/github.com/onuryilmaz/body-measurement-api/
WORKDIR /go/src/github.com/onuryilmaz/body-measurement-api/cmd
RUN GOOS=linux go build -o body-measurement

FROM ubuntu
COPY --from=builder /go/src/github.com/onuryilmaz/body-measurement-api/cmd/body-measurement /body-measurement
ENTRYPOINT ["./body-measurement"]