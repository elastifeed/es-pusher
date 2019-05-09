# Just for building
FROM golang:1.12-alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR $GOPATH/src/github.com/elastifeed/es-pusher

# Enable go Modules
ENV GO111MODULE=on

# Copy source files
COPY . .

# Fetch deps dependencies
RUN go get -d -v ./...

# Build and Install executables
RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/entrypoint.go && mkdir -p /go/bin/ && mv entrypoint /go/bin/es-pusher

# Create smallest possible docker image for production
FROM scratch

LABEL maintainer="Matthias Riegler <me@xvzf.tech>"

COPY --from=builder /go/bin/es-pusher /go/bin/es-pusher

# Entrypoint for the elasticsearch gateway
ENTRYPOINT ["/go/bin/es-pusher"]

EXPOSE 8080
