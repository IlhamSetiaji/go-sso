# Use the official Golang image as the build stage
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -ldflags "-s -w" -o main .

# Use the official Golang image as the base image for the final stage
FROM golang:1.23

# Install necessary runtime dependencies
RUN apt-get update && apt-get install -y gettext-base && rm -rf /var/lib/apt/lists/*

# Set the working directory inside the container
WORKDIR /app

# Copy the built Go application from the builder stage
COPY --from=builder /app/main .
COPY config.template.json /app/config.template.json
COPY config.json /app/config.json
COPY init-config.sh /app/init-config.sh

# Make the initialization script executable
RUN chmod +x /app/init-config.sh

# Expose the port on which the application will run
EXPOSE 3000

# Command to run the initialization script
CMD ["/app/init-config.sh"]