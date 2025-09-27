# Stage 1: Build Frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app/web

# Copy package files and install dependencies
COPY web/package.json web/package-lock.json ./
RUN npm install

# Copy the rest of the frontend source code and build
COPY web/ ./
RUN npm run build

# Stage 2: Build Backend
FROM golang:1.24-alpine AS backend-builder

WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the backend source code and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /k2ray -ldflags="-s -w" ./cmd/k2ray

# Stage 3: Final Production Image
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create a non-root user and group for security
RUN addgroup -S k2ray && adduser -S k2ray -G k2ray

# Set the working directory
WORKDIR /home/k2ray

# Copy the compiled backend binary from the backend-builder stage
COPY --from=backend-builder /k2ray .

# Copy the built frontend assets from the frontend-builder stage
COPY --from=frontend-builder /app/web/dist ./web

# Copy configuration files
COPY configs ./configs

# Change ownership of the files to the non-root user
RUN chown -R k2ray:k2ray /home/k2ray

# Switch to the non-root user
USER k2ray

# Expose the default port (assuming 8080, update if necessary)
EXPOSE 8080

# The command to run the application
CMD ["./k2ray", "--config=./configs/"]