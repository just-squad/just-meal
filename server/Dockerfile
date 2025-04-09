# syntax=docker/dockerfile:1

FROM golang:alpine AS builder
LABEL stage=gobuilder
ENV CGO_ENABLED 0
ENV GOOS linux
LABEL authors="just-squad"

# Set destination for COPY
WORKDIR /build

# Download Go modules
ADD go.mod .
ADD go.sum .
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . .

# Build
RUN go build -ldflags="-s -w" -o /app/just-meal ./cmd/just-meal/main.go

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
#EXPOSE 8080

FROM alpine
WORKDIR /app

COPY --from=builder /app/just-meal /app/just-meal
COPY --from=builder /build/cfg/appsettings.yaml /app/cfg/appsettings.yaml

# Run
CMD ["./just-meal"]