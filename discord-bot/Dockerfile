FROM golang:1.21 as builder

WORKDIR /app
COPY . .

RUN go build

FROM alpine:3.12 as runner
WORKDIR /app
COPY --from=builder /app/minecraft-server .

CMD ["./minecraft-server"]
