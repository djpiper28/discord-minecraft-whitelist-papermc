FROM golang:1.21

WORKDIR /app
COPY . .

RUN go build

ENTRYPOINT ["/app/minecraft-discord-bot"]
