# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod, go.sum
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy code
COPY backend .

# Build app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o poll-service ./cmd/api

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/poll-service .

EXPOSE 8080

CMD ["./poll-service"]
