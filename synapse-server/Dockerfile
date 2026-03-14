# ---- Build stage ----
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build the binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o synapse-server ./cmd/server

# ---- Runtime stage ----
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/synapse-server .

# Create uploads directory for media attachments
RUN mkdir -p uploads

EXPOSE 8080

CMD ["./synapse-server"]
