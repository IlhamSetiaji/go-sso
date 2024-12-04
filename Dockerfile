# Stage 1: Build the Go application
FROM golang:1.20 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main .

# Stage 2: Prepare runtime environment
FROM ubuntu:20.04

# Install dependencies (GLIBC and envsubst)
RUN apt-get update && \
    apt-get install -y ca-certificates gettext wget && \
    rm -rf /var/lib/apt/lists/*

# Upgrade GLIBC to version 2.32 if needed
RUN wget http://ftp.gnu.org/gnu/libc/glibc-2.32.tar.gz && \
    tar -xvzf glibc-2.32.tar.gz && \
    cd glibc-2.32 && \
    mkdir build && cd build && \
    ../configure --prefix=/opt/glibc-2.32 && \
    make -j$(nproc) && \
    make install && \
    rm -rf /glibc-2.32*

# Set environment for GLIBC
ENV LD_LIBRARY_PATH=/opt/glibc-2.32/lib:$LD_LIBRARY_PATH

# Set the working directory inside the container
WORKDIR /app

# Copy the built Go application from the builder stage
COPY --from=builder /app/main .
COPY config.template.json /app/config.template.json
COPY init-config.sh /app/init-config.sh

# Make the initialization script executable
RUN chmod +x /app/init-config.sh

# Expose the port on which the application will run
EXPOSE 3000

# Command to run the initialization script
CMD ["/app/init-config.sh"]
