# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 1. Copy the go.mod file first
COPY backend/go.mod ./

# 2. Copy all your source code (so Go can read your import statements)
COPY backend/. .

# 3. Download and clean up dependencies (this generates the complete go.sum)
RUN go mod download
RUN go mod tidy

# 4. Build the application binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Run stage
FROM alpine:latest

WORKDIR /app

# 5. Copy ONLY the compiled binary from the builder stage
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
