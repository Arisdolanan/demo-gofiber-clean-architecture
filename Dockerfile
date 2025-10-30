# Build stage for wkhtmltopdf
FROM surnet/alpine-wkhtmltopdf:3.22.0-0.12.6-full as wkhtmltopdf

# Build stage for Go application
FROM golang:1.23.5-alpine AS builder

# Install git and other dependencies
RUN apk add --no-cache git gcc musl-dev

# Install dependencies for wkhtmltopdf
RUN apk add --no-cache \
    libstdc++ \
    libx11 \
    libxrender \
    libxext \
    libssl3 \
    ca-certificates \
    fontconfig \
    freetype \
    ttf-dejavu \
    ttf-droid \
    ttf-freefont \
    ttf-liberation \
    build-base \
    postgresql-dev \
    gcc \
    openssl-dev \
    openssl \
    curl \
    sqlite-libs>=3.40.1-r0

# Copy wkhtmltopdf files from docker-wkhtmltopdf image
COPY --from=wkhtmltopdf /bin/wkhtmltopdf /usr/local/bin/wkhtmltopdf
COPY --from=wkhtmltopdf /bin/wkhtmltoimage /usr/local/bin/wkhtmltoimage
COPY --from=wkhtmltopdf /lib/libwkhtmltox* /usr/lib/

# Make sure binaries are executable
RUN chmod +x /usr/local/bin/wkhtmltopdf /usr/local/bin/wkhtmltoimage

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the web application with optimizations and reduced memory usage
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags '-s -w' -o main ./cmd/web

# Build the worker application with optimizations and reduced memory usage
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags '-s -w' -o worker ./cmd/worker

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Install dependencies for wkhtmltopdf in final image
RUN apk add --no-cache \
    libstdc++ \
    libx11 \
    libxrender \
    libxext \
    libssl3 \
    fontconfig \
    freetype \
    ttf-dejavu \
    ttf-droid \
    ttf-freefont \
    ttf-liberation

# Configure font cache
RUN fc-cache -fv

# Set working directory
WORKDIR /root/

# Copy the binaries from builder stage
COPY --from=builder /app/main ./main
COPY --from=builder /app/worker ./worker

# Copy config file
COPY --from=builder /app/config.json ./config.json

# Copy migrations
COPY --from=builder /app/migrations ./migrations

# Copy API documentation
COPY --from=builder /app/api ./api

# Copy wkhtmltopdf binaries and libraries
COPY --from=builder /usr/local/bin/wkhtmltopdf /usr/local/bin/wkhtmltopdf
COPY --from=builder /usr/local/bin/wkhtmltoimage /usr/local/bin/wkhtmltoimage
COPY --from=builder /usr/lib/libwkhtmltox.so* /usr/lib/

# Make sure binaries are executable in final stage
RUN chmod +x /usr/local/bin/wkhtmltopdf /usr/local/bin/wkhtmltoimage

# Create non-root user
RUN adduser -D -s /bin/sh appuser

# Change ownership of files to non-root user
RUN chown -R appuser:appuser ./

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 3000

# Command to run both services
CMD ["sh", "-c", "./worker & exec ./main"]