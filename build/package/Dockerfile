# build stage
FROM golang:1.11.1-alpine3.8 AS build
LABEL stage=intermediate

RUN apk add --no-cache git build-base
RUN go get github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/sgrzywna/milightd
COPY . /go/src/github.com/sgrzywna/milightd
RUN make

# release stage
FROM alpine:3.8
WORKDIR /app
COPY --from=build /go/src/github.com/sgrzywna/milightd/cmd/milightd/milightd /app/
ENTRYPOINT ["/app/milightd"]
