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

#### via Binary
* Create a config file for bot or use ENV variables, please check the [example](#configuration)
* Provide the db/migrations folder
* Start bot

#### via Docker
* ```docker run -it -p 3000:3000 -e DATABASE_URL="postgres_url" -e TELEGRAM_BOT_TOKEN="telegram_bot_token" to4kin/devops-school-bot:latest```

#### via AWS Lambda

Use `serverless` framework to deploy AWS Lambda. Do not forget set environment variable `AWSLAMBDA_ENABLED=true`

```yaml
service: devops-school-bot
useDotenv: true
configValidationMode: error
frameworkVersion: '>=2.61.0'

provider:
  region: "eu-central-1"
  lambdaHashingVersion: "20201221"
  name: aws
  runtime: go1.x
  logRetentionInDays: 30
  endpointType: regional
  tracing:
    apiGateway: true
    lambda: true
  ecr:
    images:
      latest:
        path: ./
  iam:
    role:
      statements:
        - Effect: "Allow"
          Resource: "*"
          Action:
            - "xray:*"

functions:
  webhook: 
    image: latest
    timeout: 15
    description: Manage students progress and provide the report if needed
    memorySize: 128
    environment:
      AWSLAMBDA_ENABLED: true
      DATABASE_URL: ${env:DATABASE_URL}
      TELEGRAM_BOT_TOKEN: ${env:TELEGRAM_BOT_TOKEN}
      
    events:
      - http:
          path: /webhook
          method: POST
          cors: false

```

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

### Configuration

Availiable file formats: `JSON, TOML, YAML, HCL, envfile and Java properties`

example-config.toml

```toml
log_level = "debug"
debug_mode = false

[apiserver]
bind_addr = ":3000"

[apiserver.cron]
enabled = true
fullreport = true
schedule = "0 15 * * FRI"

[awslambda]
enabled = false

[database]
url = "postgres://localhost/devops_school?user=postgres&password=example&sslmode=disable"
migrations = "db/migrations"

[telegram_bot]
token = "TEST_TELEGRAM_TOKEN"
verbose = false

```

example-config.yaml:

```yaml
log_level: debug
debug_mode: false

apiserver:
  bind_addr: :3000

  cron:
    enabled: true
    fullreport: true
    schedule: "0 15 * * FRI"

awslambda:
  enabled: false

database:
  url: postgres://localhost/devops_school?user=postgres&password=example&sslmode=disable
  migrations: db/migrations

telegram_bot:
  token: TEST_TELEGRAM_TOKEN
  verbose: false

```

Or, you can use ENV variables to update config:
```bash
LOG_LEVEL="debug"
DEBUG_MODE="false"

APISERVER_BIND_ADDR=":3000"
APISERVER_CRON_ENABLED="true"
APISERVER_CRON_FULLREPORT="true"
APISERVER_CRON_SCHEDULE="0 15 * * FRI"

AWSLAMBDA_ENABLED="false"

DATABASE_URL="postgres://localhost/devops_school?user=postgres&password=example&sslmode=disable"
DATABASE_MIGRATIONS="db/migrations"

TELEGRAM_BOT_TOKEN="TEST_TELEGRAM_TOKEN"
TELEGRAM_BOT_VERBOSE="false"

```

### Cron job

By default, the Bot will send a full report each Friday at 15:00 UTC time to all active school chats. It's the same report as /fullreport command. If you would like to use short report (as /report command) - set `fullreport = false`

## Docker Compose example

```
version: '3.7'

services:
  postgres:
    image: postgres:13.4
    container_name: postgres
    restart: always
    environment:
      POSTGRES_DB: devops_school
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: example
    volumes:
      - db_data:/var/lib/postgresql/data
  application:
    image: to4kin/devops-school-bot:latest
    container_name: devops-school-bot
    restart: always
    environment:
      DATABASE_URL: postgres://postgres/devops_school?user=postgres&password=example&sslmode=disable
      TELEGRAM_BOT_TOKEN: TELEGRAM_BOT_TOKEN
    depends_on:
      - postgres

volumes:
  db_data:

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

## Debug mode

When debug_mode is enabled the bot will provide additional http endpoint /debug/pprof for profiling data in the format expected by the pprof visualization tool

## Student types

There're 2 types of students: full course and module course.

We have a dynamic homework list, which means that at the beginning of the school the list is empty and populated when the homework is provided by the student. This means that the student should provide all homework which was provided by other students.

On the other hand, if the student was joined for module course, the list of homework depends on the modules. The module list is also populated by students.

We need to follow one simple instruction: the homework hashtag should be #<MODULE_NAME><SEQUENTIAL_NUMBER>, for example `#cicd1 #cicd2 #cloud1 #ansible1`
