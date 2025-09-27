# Stage 1: Build the Go application
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies first, leveraging Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application, creating a statically linked binary
# This is important for running in a 'distroless' container
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /k2ray cmd/k2ray/main.go

# Stage 2: Create the final, minimal image
FROM gcr.io/distroless/static-debian12

# Set the working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /k2ray /app/k2ray

# Copy the configuration files needed by the application
COPY configs/ /app/configs/

# Expose the port the application runs on
EXPOSE 8080

# Set the user to non-root for security
USER nonroot:nonroot

# The command to run the application
ENTRYPOINT ["/app/k2ray"]