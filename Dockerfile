FROM alpine:3.14.2
COPY bin/devops-school-bot.linux.amd64 /app/devops-school-bot
ENTRYPOINT ["/app/devops-school-bot", "start"]
