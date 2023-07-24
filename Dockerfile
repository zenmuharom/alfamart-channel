# ================ Start build
FROM golang:1.19.2-alpine3.16 AS builder
WORKDIR /app
COPY . .
# installing gcc
RUN apk --no-cache add make gcc libtool musl-dev ca-certificates dumb-init
RUN CGO_ENABLED=1 go build -ldflags="-s -w" -o main main.go
# ================ End build

# ================ Start running app
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
COPY example.env app.env
COPY start.sh .
EXPOSE 80
CMD ["/app/main"]
# ================ End running app