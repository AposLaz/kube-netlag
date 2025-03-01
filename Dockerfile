FROM golang:1.24 AS builder

RUN apt-get update && apt-get install -y git && \
    useradd -m -s /bin/bash gouser

# Set working directory
WORKDIR /app

RUN chown -R gouser:gouser /app

COPY pkg/go.mod pkg/go.sum ./

# Download Go module dependencies
RUN go mod download

COPY ./pkg ./pkg

USER gouser
WORKDIR /app/pkg
RUN go build -o /app/kube-netlag

FROM ubuntu:24.10

# Install dependencies for netperf
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates netperf \
    && rm -rf /var/lib/apt/lists/* && \
    useradd -m -s /bin/bash gouser

WORKDIR /app

COPY --from=builder /app/kube-netlag .

RUN chown -R gouser:gouser /app
USER gouser

CMD ["./kube-netlag"]
