ARG GO_VERSION=1.22
FROM golang:${GO_VERSION}-alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod tidy

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o app ./cmd

WORKDIR /bin

# Copy binary from build to main folder
RUN cp /build/app .

FROM golang:${GO_VERSION}-alpine
COPY --from=builder /bin /
COPY --from=builder /build/docs ./docs
CMD ["/app"]