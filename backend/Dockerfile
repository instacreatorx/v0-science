# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# ✨ FIX: Copy from the backend directory explicitly
COPY backend/go.mod ./

# Fetch dependencies and generate go.sum
RUN go mod download
RUN go mod tidy

# ✨ FIX: Copy the rest of the backend source code
COPY backend/. .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .
