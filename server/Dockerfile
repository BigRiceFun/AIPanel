FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY server/go.mod server/go.sum* ./
RUN go mod download

COPY server/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/aipanel-server .

FROM alpine:3.20

WORKDIR /app
RUN apk add --no-cache bash ca-certificates tzdata

COPY --from=builder /bin/aipanel-server /usr/local/bin/aipanel-server
COPY config.yaml ./config.yaml

EXPOSE 8080
CMD ["aipanel-server"]
