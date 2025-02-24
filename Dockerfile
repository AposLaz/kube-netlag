# ---- Builder Stage (Go) ----
FROM golang:1.24 AS builder

# Install git
RUN apt-get update && apt-get install -y git

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum for dependency caching
COPY src/go.mod src/go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy source code
COPY ./src ./src

# Change to src directory and build
WORKDIR /app/src
RUN go build -o /app/kube-netlag

# ---- Final Stage ----
FROM ubuntu:24.10

# Install dependencies for netperf and bind-tools
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    linux-headers-generic \
    bind9-dnsutils \
    ca-certificates \
    netperf \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy the built Go binary from the builder stage
COPY --from=builder /app/kube-netlag .

# Run the Go application
CMD ["./kube-netlag"]
