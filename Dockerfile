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
RUN go build -o main .

# Use the official Alpine image as the base image for the final stage
FROM alpine:3.18

# Set timezone ke non-interactive dan install dependencies
RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates \
    wget \
    gcc \
    make \
    gawk \
    bison \
    python3 \
    gettext \
    build-essential \
    tzdata && \
    ln -fs /usr/share/zoneinfo/Etc/UTC /etc/localtime && \
    echo "Etc/UTC" > /etc/timezone && \
    dpkg-reconfigure -f noninteractive tzdata && \
    rm -rf /var/lib/apt/lists/*

# Install GLIBC 2.32
RUN wget http://ftp.gnu.org/gnu/libc/glibc-2.32.tar.gz && \
    tar -xvzf glibc-2.32.tar.gz && \
    cd glibc-2.32 && \
    mkdir build && cd build && \
    ../configure --prefix=/opt/glibc-2.32 && \
    make -j$(nproc) && \
    make install && \
    rm -rf /glibc-2.32*




# Set the working directory inside the container
WORKDIR /app

# Copy the built Go application from the builder stage
# Copy the built Go application from the builder stage
COPY --from=builder /app/main .
COPY config.template.json /app/config.template.json
COPY init-config.sh /app/init-config.sh
# Convert init-config.sh to LF format and make it executable
RUN chmod +x /app/init-config.sh

# Expose the port on which the application will run
EXPOSE 3000

# Command to run the initialization script
CMD ["/app/init-config.sh"]
