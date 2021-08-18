# DevOps School Bot

Manage students progress and provide the report if needed

## HowTo

### Usage

```bash
DevOps School Bot manage students progress and provide the report

Usage:
  devops-school-bot [command]

Available Commands:
  help        Help about any command
  start       Start DevOps School Bot
  version     Print version

Flags:
  -h, --help   help for devops-school-bot

Use "devops-school-bot [command] --help" for more information about a command.
```

### Start bot

```bash
Start DevOps School Bot with config file
Simply execute devops-school-bot start -c path/to/config/file.toml
or skip this flag to use default path

Usage:
  devops-school-bot start [flags]

Flags:
  -c, --config-path string   path to config file (default "configs/devopsschoolbot.toml")
  -h, --help                 help for start
```

### Config file

```toml
bind_addr = ":3000"
log_level = "debug"
database_url = "postgres://localhost/devops_school_dev?user=postgres&password=example&sslmode=disable"

[telegram_bot]
token = "1949550059:AAHTvp0Zm5ABVDKL8LVHAYkS-PEEGGZnEJE"
verbose = false
```

## Bot commands

There're two types of telegram chat where the bot is accepting commands: `private` and `group`

### Private chat

| Command | Roles | Description |
| ------- | ----- | ----------- |
| /schools | **Supersuser** | Provide the interface to manage schools |
| /users | **Supersuser** | Provide the interface to manage users |
| /start  | User | Add user to the database for future use |
| /help   | User | Help message | 

### Group chat

| Command | Roles | Description |
| ------- | ----- | ----------- |
| /start  | **Supersuser** | Start school with name == title |
| /finish | **Supersuser** | Finish school. Homeworks are not accepted after school finished |
| /report | **Supersuser** | Provide school progress |
| /bigreport | **Supersuser** | Provide school progress with homework list |
| /join   | User | Add user to school as student |
| /myreport | User | Provide school progress |
| /homeworks | User | Provide the homework list |
| /help   | User | Help message | 

To add homework, use the work hashtag along with the **#homework**, for example:

**#homework #cicd**

### For BotFather

```
start - Group: Start school (Superuser). Private: Add user to database
finish - Group: Finish school (Superuser)
join - Group: Add user to school as student
myreport - Group: School progress
report - Group: School progress (Superuser)
bigreport - Group: School progress with homework list (Superuser)
homeworks - Group: homework list
schools - Private: Interface to manage schools (Superuser)
users - Private: Interface to manage users (Superuser)
help - All: Help message
```