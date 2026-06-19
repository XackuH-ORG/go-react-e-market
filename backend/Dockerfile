FROM golang:1.26.3-alpine AS builder

WORKDIR /build

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/backend ./cmd/

# Final stage
FROM alpine:latest

# Add certificates for external requests and create non-root user
RUN apk --no-cache add ca-certificates && \
    adduser -D -g '' appuser

WORKDIR /app

# Copy binary from builder
COPY --from=builder --chown=appuser:appuser /app/backend .

# Copy config directory
COPY --chown=appuser:appuser config/ config/

USER appuser

# Expose port (default 8080)
EXPOSE 8080

CMD ["./backend"]
