FROM golang:alpine3.17 as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o api-url-shortener cmd/main.go

FROM golang:alpine3.17
WORKDIR /app
COPY --from=builder /app/api-url-shortener .
COPY --from=builder /app/.env .
CMD ["./api-url-shortener"]