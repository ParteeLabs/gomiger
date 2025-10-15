# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git for go modules
RUN apk add --no-cache git

# Copy go.work and module files
COPY go.work go.work.sum ./
COPY core/go.mod core/go.sum ./core/
COPY mongomiger/go.mod mongomiger/go.sum ./mongomiger/
COPY examples/go.mod examples/go.sum ./examples/

# Download dependencies
RUN cd core && go mod download
RUN cd mongomiger && go mod download
RUN cd examples && go mod download

# Copy source code
COPY . .

# Build the gomiger-init binary
RUN cd core/cmd/gomiger-init && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /gomiger-init .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /gomiger-init .

# Make the binary executable
RUN chmod +x ./gomiger-init

# Create a directory for migrations
RUN mkdir -p /app/migrations

WORKDIR /app

# Set the entrypoint
ENTRYPOINT ["/root/gomiger-init"]