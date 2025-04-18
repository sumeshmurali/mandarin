# Stage 1: Build the Go binary
FROM golang:1.24 AS builder

RUN mkdir /mandarin

WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app/ 

# Stage 2: Create the final image
FROM alpine:latest

WORKDIR /app

# Copy the Go binary from the builder stage
COPY --from=builder /app/app .

# Set the command line parameter for configuration
# CMD ["./app", "-config", "config.json"]
ENV ENV_DOCKER=true
ENTRYPOINT [ "./app" ]