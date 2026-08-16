FROM golang:1.25-alpine AS builder

WORKDIR /app

# Cache dependencies separately from source code — this layer only
# rebuilds when go.mod/go.sum change, not on every code edit
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# ---- Runtime stage ----
FROM alpine:3.20

WORKDIR /app

# Certificates needed for outbound HTTPS (e.g. if MongoDB Atlas is ever used)
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]