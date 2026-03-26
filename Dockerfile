FROM golang:1.24 AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o tk-api .

FROM debian:stable-slim

WORKDIR /app
COPY --from=builder /app/tk-api .
COPY etc/tk-api.yaml ./etc/tk-api.yaml

EXPOSE 8088

CMD ["./tk-api","-f","etc/tk-api.yaml"]