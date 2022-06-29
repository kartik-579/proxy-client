FROM golang:1.18.3-alpine3.16 AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy the code into the container
COPY ./main.go .

# Build the application
RUN go build -o proxy-client ./main.go

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/proxy-client .

EXPOSE 8080

# Build a small image
FROM alpine:3.16.0
#COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /dist/proxy-client /
ENTRYPOINT ["/proxy-client"]