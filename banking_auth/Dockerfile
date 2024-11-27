# Use Golang official image
FROM golang:1.23

# Set environment variables
ENV SERVER_ADDRESS="0.0.0.0" \
    SERVER_PORT="8181" \
    MONGODB_URI="mongodb://localhost:27017"

WORKDIR /app-auth

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy application code into container
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the application port
EXPOSE 8181

# Run the application
CMD ["./main"]
