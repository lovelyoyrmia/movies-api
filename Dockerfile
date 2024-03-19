# Build stage
FROM golang:1.21-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o main cmd/main.go

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY wait-for.sh .
COPY app.env .
COPY docs ./docs
COPY pkg ./pkg
COPY migrations ./migrations

EXPOSE 50054
CMD [ "/app/main" ]