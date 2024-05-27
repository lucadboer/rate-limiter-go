FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o rate-limiter .

FROM alpine

WORKDIR /app

COPY --from=builder /app/rate-limiter .
COPY .env .   

ENTRYPOINT ["./rate-limiter"]
