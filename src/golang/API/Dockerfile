FROM golang:alpine AS builder

# Cross compile, either supplied as a --build-arg or automagically by docker buildx
ARG TARGETPLATFORM

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARC=$TARGETPLATFORM \
    GOBIN=$GOPATH/bin

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY * /build/

# Build the application
RUN go build

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/Sabatical ./main

# Build a small image
FROM scratch

COPY --from=builder /dist/main /main

# Command to run
ENTRYPOINT ["/main"]
