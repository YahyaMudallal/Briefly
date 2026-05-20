# =================================
# Build stage
# ================================
FROM golang:1.22-alpine AS builder

# Install git and ca-certificates to allow fetching dependencies
RUN apk add --no-cache git ca-certificates

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to the working directory
COPY go.mod go.sum* ./

# Download the dependencies
RUN go mod download

# Copy the entire source code to the working directory
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o news-website-backend

# =================================
# Run stage
# ================================
FROM alpine:latest

# Install ca-certificates and tzdata for HTTPS requests and timezone support
RUN apk add --no-cache ca-certificates tzdata

# Configure the timezone
ENV TZ=Europe/Paris

# Set the working directory inside the container
WORKDIR /app

# Get the binary from the builder stage
COPY --from=builder /app/news-website-backend .

# Expose the port
EXPOSE 10000

# Launch the application
CMD ["./news-website-backend"]
