# Use the official Golang image as the build stage
FROM golang:1.23-alpine AS builder

# Install necessary build dependencies
RUN apk add --no-cache gcc musl-dev

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

# Use the official Alpine image as the base image for the final stage
FROM alpine:3.18
RUN apk add --no-cache gettext

# Set the working directory inside the container
WORKDIR /app

# Copy the built Go application from the builder stage
# Copy the built Go application from the builder stage
COPY --from=builder /app/main .
COPY config.template.json /app/config.template.json
COPY init-config.sh /app/init-config.sh
COPY config.json /app/config.json
# Convert init-config.sh to LF format and make it executable
RUN chmod +x /app/init-config.sh /app/main

# Expose the port on which the application will run
EXPOSE 3000

# Command to run the initialization script
CMD ["/app/main"]
