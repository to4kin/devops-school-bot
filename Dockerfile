# Build stage
FROM artifactory.devops.telekom.de/hub.docker.com/golang:1.17 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make all

# CMD stage
FROM artifactory.devops.telekom.de/hub.docker.com/alpine:3.14.2
COPY --from=builder /app/bin/gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot /app/devops-school-bot
ENTRYPOINT ["/app/devops-school-bot", "start"]
