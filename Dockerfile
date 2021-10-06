FROM alpine:3.14.2
COPY db/ /app/db
COPY configs/ /app/configs
COPY bin/devops-school-bot.linux.amd64 /app/devops-school-bot
WORKDIR /app
ENTRYPOINT ["./devops-school-bot", "start"]
