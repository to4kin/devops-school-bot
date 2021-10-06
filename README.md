[![ci-build](https://github.com/to4kin/devops-school-bot/actions/workflows/ci-build.yml/badge.svg?branch=master)](https://github.com/to4kin/devops-school-bot/actions/workflows/ci-build.yml)
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/to4kin/devops-school-bot.svg)](https://github.com/to4kin/devops-school-bot)
[![GitHub release](https://img.shields.io/github/release/to4kin/devops-school-bot.svg)](https://GitHub.com/to4kin/devops-school-bot/releases/)
[![codecov](https://codecov.io/gh/to4kin/devops-school-bot/branch/master/graph/badge.svg?token=FHZ0TRMG92)](https://codecov.io/gh/to4kin/devops-school-bot)

# DevOps School Bot

Manage students progress and provide the report if needed

## Precondition

* PostgreSQL 13
## HowTo

### Start

* Add bot to the chat
* Make bot admin (to read all messages, not only commands)
* Set a webhook to bot https address: `curl https://api.telegram.org/bot<BOT_TOKEN>/setWebhook\?url\=<BOT_HTTPS_URL>`
* Create a config file for bot or use ENV variables, please check the example
* Provide the db/migrations folder (or use docker image)
* Start bot

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

config.toml

```toml
bind_addr = ":3000"
log_level = "debug"

[database]
url = "postgres://localhost/devops_school_dev?user=postgres&password=example&sslmode=disable"
migrations = "db/migrations"

[telegram_bot]
token = "TEST_TELEGRAM_TOKEN"
verbose = false
```

config.yaml:

```yaml
bind_addr: :3000
log_level: debug

database:
  url: postgres://localhost/devops_school_dev?user=postgres&password=example&sslmode=disable
  migrations: db/migrations

telegram_bot:
  token: TEST_TELEGRAM_TOKEN
  verbose: false
```

Or, you can use ENV variables to update config:
```bash
BIND_ADDR=":3000"
LOG_LEVEL="debug"
DATABASE_URL="postgres://localhost/devops_school_dev?user=postgres&password=example&sslmode=disable"
DATABASE_MIGRATIONS="db/migrations"
TELEGRAM_BOT_TOKEN="TEST_TELEGRAM_TOKEN"
TELEGRAM_BOT_VERBOSE="false"
```

## Docker

```
version: '3.7'

services:
  postgres:
    image: postgres:13.4
    container_name: postgres
    restart: always
    environment:
      POSTGRES_DB: devops_school
      POSTGRES_PASSWORD: strongpassword
    volumes:
      - postgresql_data:/var/lib/postgresql/data
  application:
    image: to4kin/devops-school-bot:latest
    container_name: devops-school-bot
    restart: always
    environment:
      DATABASE_URL: postgres://postgres/devops_school_dev?user=postgres&password=strongpassword&sslmode=disable
    depends_on:
      - postgres
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

## Student types

There're 2 types of students: full course and module course.

We have a dynamic homework list, which means that at the beginning of the school the list is empty and populated when the homework is provided by the student. This means that the student should provide all homework which was provided by other students.

On the other hand, if the student was joined for module course, the list of homework depends on the modules. The module list is also populated by students.

We need to follow one simple instruction: the homework hashtag should be #<MODULE_NAME><SEQUENTIAL_NUMBER>, for example `#cicd1 #cicd2 #cloud1 #ansible1`