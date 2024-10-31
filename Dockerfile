# -----------------------------
# Stage 1: Build the Go binary
# -----------------------------
    FROM golang:1.23-alpine AS builder

    # Install git for fetching dependencies
    RUN apk add --no-cache git
    
    # Set the working directory inside the container
    WORKDIR /app
    
    # Copy go.mod and go.sum files
    COPY go.mod go.sum ./
    
    # Download dependencies
    RUN go mod download
    
    # Copy the source code
    COPY . .
    
    # Build the Go app
    RUN go build -o ssh_key_scanner
    
    # -----------------------------
    # Stage 2: Create a minimal image
    # -----------------------------
    FROM alpine:latest
    
    # Install CA certificates (if needed)
    RUN apk add --no-cache ca-certificates
    
    # Set the working directory
    WORKDIR /root/
    
    # Copy the compiled binary from the builder stage
    COPY --from=builder /app/ssh_key_scanner .
    
    # Set the entrypoint
    ENTRYPOINT ["./ssh_key_scanner"]
    