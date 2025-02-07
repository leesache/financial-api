# docker/app/Dockerfile
# Stage 1: Build the Go application
FROM golang:1.23.3-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o financial-api ./main.go

# Stage 2: Create a minimal runtime image
FROM alpine:latest
RUN apk add --no-cache bash
WORKDIR /root/
COPY --from=builder /app/financial-api .
EXPOSE 8080
CMD ["./financial-api"]