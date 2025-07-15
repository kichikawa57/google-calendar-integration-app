# Development stage
FROM golang:1.24-alpine

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum first (for caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code, including main.go
COPY . .

# Expose port (optional)
EXPOSE 8080

# Run the app
CMD ["go", "run", "main.go"]