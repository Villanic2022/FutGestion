# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies.
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o futbol-api .

# Run stage
FROM alpine:latest

WORKDIR /app/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/futbol-api .
# If you have other files like static assets or templates, copy them too
COPY --from=builder /app/logo.png .

# Expose port (adjust if your API uses a different one)
EXPOSE 8080

# Command to run the executable
CMD ["./futbol-api"]
