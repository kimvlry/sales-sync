FROM golang:1.24.5-alpine AS builder

WORKDIR /app

COPY go.work ./
COPY services/user-service/go.mod services/user-service/go.sum ./services/user-service/
COPY shared/go.mod shared/go.sum ./shared/

RUN go work sync
RUN go work use ./services/user-service ./shared
RUN go mod download

COPY . .

WORKDIR /app/services/user-service
RUN go build -o user-service ./cmd/main.go

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/services/user-service/user-service .

EXPOSE 8081

CMD ["./user-service"]
