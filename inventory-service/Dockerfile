# Build stage
FROM golang:1.24.9-alpine3.22 AS builder

# Declare the current working directory inside the base image
WORKDIR /app

# Copy all necessary files from where we run the docker build command
# to the working directory
COPY . .

# Bundle app source
RUN go build -o main main.go

RUN GRPC_HEALTH_PROBE_VERSION=v0.4.5 && \
    wget -qO/app/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /app/grpc_health_probe

# Run stage
FROM scratch

WORKDIR /app

# Bundle executable app and (files) dependencies
COPY --from=builder /app/main .
COPY --from=builder /app/grpc_health_probe /bin/grpc_health_probe

COPY app.yaml.dist ./app.yaml
COPY db/migration ./db/migration

# EXPOSE act as a documentation about which ports are intended to be published
EXPOSE 9090/tcp

CMD [ "/app/main" ]
