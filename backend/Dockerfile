# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 1. 🛑 REMOVED go.sum from the explicit copy step
COPY go.mod ./

# 2. ✨ ADDED go mod download before tidy to fetch dependencies
RUN go mod download
RUN go mod tidy

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Run stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
