FROM golang:1.19 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOARCH=arm64 go build -o mycli main.go

FROM debian:bullseye-slim
WORKDIR /root/
COPY --from=builder /app/mycli .
ENTRYPOINT ["./mycli"]
CMD ["--help"]