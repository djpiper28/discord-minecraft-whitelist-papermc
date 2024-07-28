FROM golang:1.21

WORKDIR /app
COPY . .

ENV DISCORD_TOKEN=change_me
ENV DATABASE_RUl=change_me
RUN go build

ENTRYPOINT ["/app/minecraft-discord-bot"]
