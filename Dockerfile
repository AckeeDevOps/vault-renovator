FROM golang:1.11.1-alpine3.8 AS builder
MAINTAINER stepan.vrany@ackee.cz
WORKDIR /go/src/github.com/AckeeDevOps/vault-renovator
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/AckeeDevOps/vault-renovator/app .
CMD ["./app"]
