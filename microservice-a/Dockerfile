# Use official Golang image
FROM golang:1.21-alpine

# Set working directory
WORKDIR /app

# Copy go.mod and download dependencies
COPY go.mod ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN go build -o microservice-a main.go

# Expose port
EXPOSE 8080

# Run the service
CMD ["./microservice-a"]
