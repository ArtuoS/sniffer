FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /sniffer ./cmd/sniffer

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /sniffer /sniffer
EXPOSE 8081
ENTRYPOINT ["/sniffer"]
