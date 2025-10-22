# Build stage
FROM golang:1.24.9-alpine3.22 AS builder

# Declare the current working directory inside the base image
WORKDIR /app

# Copy all necessary files from where we run the docker build command
# to the working directory
COPY . .

# Bundle app source
RUN go build -mod=vendor -o main main.go

# Run stage
FROM scratch

WORKDIR /app

# Bundle executable app and (files) dependencies
COPY --from=builder /app/main .

COPY app.env.example ./app.env
COPY db/migration ./db/migration

# EXPOSE act as a documentation about which ports are intended to be published
EXPOSE 8000/tcp

CMD [ "/app/main" ]
