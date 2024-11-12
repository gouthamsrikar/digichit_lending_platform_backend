# Use an official Go image as a base
FROM golang:1.23-alpine

# Install build dependencies for SQLite
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Enable CGO
ENV CGO_ENABLED=1

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o my-go-app

# Expose port 80 to access the app
EXPOSE 80

# Run the executable
CMD ["./my-go-app"]
