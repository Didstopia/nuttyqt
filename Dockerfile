# syntax=docker/dockerfile:1

FROM golang:1.19-alpine AS builder
WORKDIR /app
RUN go install github.com/cosmtrek/air@latest
CMD ["air", "-c", ".air.toml"]
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o nuttyqt .

FROM scratch
WORKDIR /app
COPY --from=builder /app/nuttyqt .
CMD ["/app/nuttyqt"]
