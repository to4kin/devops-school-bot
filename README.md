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

[database]
url = "postgres://localhost/devops_school_dev?user=postgres&password=example&sslmode=disable"
migrations = "db/migrations"

[telegram_bot]
token = "1949550059:AAHTvp0Zm5ABVDKL8LVHAYkS-PEEGGZnEJE"
verbose = false
```

## Bot commands

There're two types of telegram chat where the bot is accepting commands: `private chat` and `group chat`

##### User Commands
**/start** - Add user to database
**/joinstudent** - Join school as student
**/joinmodule** - Join school as listener
**/myreport** - Your progress
**/homeworks** - Homeworks list
**/help** - Help message

##### Superuser Commands
**/schools** - Manage schools
**/startschool** - Start school
**/stopschool** - Finish school
**/report** - School progress
**/fullreport** - School progress with homework list

**/users** - Manage users
**/setsuperuser** - Set Superuser
**/unsetsuperuser** - Unset Superuser

To add homework, use the work hashtag along with the `#homework`, for example: `#homework #cicd`

#### For BotFather

```
start - Add user to database
joinstudent - Join school as student
joinmodule - Join school as listener
myreport - Your progress
homeworks - Homeworks list
help - Help message
schools - Manage schools
startschool - Start school
stopschool - Finish school
report - School progress
fullreport - School progress with homework list
users - Manage users
setsuperuser - Set Superuser
unsetsuperuser - Unset Superuser
```