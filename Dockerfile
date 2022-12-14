FROM golang:1.18-alpine3.15 as Builder
COPY . /go/src/app
WORKDIR /go/src/app
RUN go mod tidy -compat=1.18
RUN go build -o /go/src/receiver/bin/receiver /go/src/receiver/cmd/receiver/receiver.go

FROM alpine:3.15
COPY --from=Builder /go/src/app/bin/* /go/src/app/