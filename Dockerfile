# OCPP 1.6 example

# =================================================
# Build container
# =================================================
FROM golang:1.18-alpine3.15 AS builder

ENV GO111MODULE on
# Create folder for project
WORKDIR $GOPATH/src/github.com/CoderSergiy/ocpp16-go
# Copy branch to the container
COPY . .
# Fetch dependencies.
RUN go get -v -t ./...
# Build the binary.
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/server server.go


# =================================================
# Launch server on the small container
# =================================================
FROM alpine

COPY --from=builder /go/bin/server /bin/server

# Add CA certificates
# It currently throws a warning on alpine: WARNING: ca-certificates.crt does not contain exactly one certificate or CRL: skipping.
# Ignore the warning.
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/* && update-ca-certificates
# Copy configs file to the container
COPY example/configs.json /tmp
# Create folder for the log files in the container
RUN mkdir /tmp/logs
# Since running as a non-root user, port bindings < 1024 is not possible
# 8000 for HTTP; 8443 for HTTPS;
EXPOSE 8080

CMD ["/bin/server"]