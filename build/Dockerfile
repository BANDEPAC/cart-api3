FROM golang:1.22.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/cart-api ./cmd/cart-api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/cart-api .

COPY migrations ./migrations

COPY app.env .

CMD ["./cart-api"]