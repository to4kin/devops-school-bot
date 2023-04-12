provider "aws" {
  region = local.region

  # Make it faster by skipping something
  skip_metadata_api_check     = true
  skip_region_validation      = true
  skip_credentials_validation = true
  skip_requesting_account_id  = true
}

data "aws_caller_identity" "current" {}
data "aws_region" "current" {}
data "aws_availability_zones" "available" {}
data "aws_ecr_authorization_token" "token" {}

locals {
  name   = "devops-school-bot"
  region = "eu-central-1"

  vpc_cidr = "10.0.0.0/16"
  azs      = slice(data.aws_availability_zones.available.names, 0, 3)
}

################################################################################
# AWS Lambda
################################################################################

provider "docker" {
  registry_auth {
    address  = format("%v.dkr.ecr.%v.amazonaws.com", data.aws_caller_identity.current.account_id, local.region)
    username = data.aws_ecr_authorization_token.token.user_name
    password = data.aws_ecr_authorization_token.token.password
  }
}

module "lambda_function" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "~> 4.12"

  function_name = local.name
  runtime       = "go1.x"
  memory_size   = 128
  timeout       = 60
  environment_variables = {
    AWSLAMBDA_ENABLED  = true
    DATABASE_URL       = format("postgres://%v/%v?user=%v&password=%v", module.db.db_instance_endpoint, module.db.db_instance_name, module.db.db_instance_username, random_password.master_password.result)
    TELEGRAM_BOT_TOKEN = var.telegram_bot_token
  }

  create_package = false

  policy_json = <<EOF
  {
      "Version": "2012-10-17",
      "Statement": [
          {
              "Effect": "Allow",
              "Action": [
                  "xray:*"
              ],
              "Resource": ["*"]
          }
      ]
  }
  EOF

  allowed_triggers = {
    APIGatewayAny = {
      service    = "apigateway"
      source_arn = "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${module.api_gateway.api_gateway_id}/*/*/${local.name}"
    }
  }

  publish = true

  ##################
  # Container Image
  ##################
  image_uri     = module.docker_image.image_uri
  package_type  = "Image"
  architectures = ["x86_64"]

  vpc_subnet_ids         = module.vpc.private_subnets
  vpc_security_group_ids = [module.lambda_security_group.security_group_id]
  attach_network_policy  = true

  tags = merge(
    var.additional_tags,
    var.security_class_tags
  )
}

module "docker_image" {
  source = "./modules/docker-build"

  create_ecr_repo = true
  ecr_repo        = local.name

  image_tag   = "latest"
  source_path = "${path.cwd}/../"
  platform    = "linux/amd64"

  ecr_repo_tags = merge(
    var.additional_tags,
    var.security_class_tags
  )

  triggers = {
    bin_sha1    = sha1("${path.cwd}/../bin/devops-school-bot.linux.amd64")
    config_sha1 = sha1(join("", [for f in fileset(path.module, "${path.cwd}/../configs") : filesha1(f)]))
    db_sha1     = sha1(join("", [for f in fileset(path.module, "${path.cwd}/../db/migrations") : filesha1(f)]))
  }
}

################################################################################
# API Gateway
################################################################################

module "api_gateway" {
  source = "./modules/api-gateway"

  api_gateway_name            = local.name
  api_gateway_http_method     = "POST"
  api_gateway_integration_uri = module.lambda_function.lambda_function_invoke_arn
  api_gateway_stage_name      = "webhook"

  api_gateway_tags = merge(
    var.additional_tags,
    var.security_class_tags,
    {
      "dtit:sec:NetworkLayer" = "application"
    }
  )
}


################################################################################
# RDS Module
################################################################################

resource "random_password" "master_password" {
  length  = 16
  special = false
}

resource "aws_secretsmanager_secret" "rds_credentials" {
  name = local.name
}

resource "aws_secretsmanager_secret_version" "rds_credentials" {
  secret_id     = aws_secretsmanager_secret.rds_credentials.id
  secret_string = <<EOF
{
  "username": "${module.db.db_instance_username}",
  "password": "${random_password.master_password.result}",
  "engine": "${module.db.db_instance_engine_version_actual}",
  "host": "${module.db.db_instance_address}",
  "port": ${module.db.db_instance_port},
}
EOF
}

module "db" {
  source  = "terraform-aws-modules/rds/aws"
  version = "~> 5.6"

  identifier = local.name

  # All available versions: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_PostgreSQL.html#PostgreSQL.Concepts
  engine               = "postgres"
  engine_version       = "14"
  family               = "postgres14" # DB parameter group
  major_engine_version = "14"         # DB option group
  instance_class       = "db.t3.medium"

  allocated_storage     = 20
  max_allocated_storage = 50

  # NOTE: Do NOT use 'user' as the value for 'username' as it throws:
  # "Error creating DB Instance: InvalidParameterValue: MasterUsername
  # user cannot be used as it is a reserved word used by the engine"
  db_name  = "devops_school"
  username = "devops_school"
  password = random_password.master_password.result
  port     = 5432

  multi_az               = false
  publicly_accessible    = false
  storage_encrypted      = true
  db_subnet_group_name   = module.vpc.database_subnet_group
  vpc_security_group_ids = [module.rds_security_group.security_group_id]

  maintenance_window = "Mon:00:00-Mon:03:00"
  backup_window      = "03:00-06:00"

  enabled_cloudwatch_logs_exports = ["postgresql", "upgrade"]
  create_cloudwatch_log_group     = true

  backup_retention_period = 1
  copy_tags_to_snapshot   = true
  skip_final_snapshot     = true
  deletion_protection     = false

  performance_insights_enabled          = true
  performance_insights_retention_period = 7
  create_monitoring_role                = true
  monitoring_interval                   = 60
  monitoring_role_name                  = "rds-monitoring-role-name"
  monitoring_role_use_name_prefix       = true
  monitoring_role_description           = "RDS monitoring role"
  parameters = [
    {
      name  = "autovacuum"
      value = 1
    },
    {
      name  = "client_encoding"
      value = "utf8"
    },
    {
      name  = "rds.force_ssl"
      value = 1
    }
  ]

  tags = merge(
    var.additional_tags,
    var.security_class_tags,
    var.schedule_tags
  )

  db_option_group_tags = merge(
    var.additional_tags,
    var.security_class_tags,
    {
      "Sensitive" = "low"
    }
  )

  db_parameter_group_tags = merge(
    var.additional_tags,
    var.security_class_tags,
    {
      "Sensitive" = "low"
    }
  )
}

################################################################################
# Supporting Resources
################################################################################

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 3.19"

  name = "${local.name}-vpc"
  cidr = local.vpc_cidr

  azs              = local.azs
  public_subnets   = [for k, v in local.azs : cidrsubnet(local.vpc_cidr, 8, k)]
  private_subnets  = [for k, v in local.azs : cidrsubnet(local.vpc_cidr, 8, k + 3)]
  database_subnets = [for k, v in local.azs : cidrsubnet(local.vpc_cidr, 8, k + 6)]

  enable_nat_gateway           = true
  single_nat_gateway           = true
  create_database_subnet_group = true

  tags = merge(
    var.additional_tags,
    var.security_class_tags
  )
}

module "rds_security_group" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 4.0"

  name        = "${local.name}-rds-security-group"
  description = "Complete PostgreSQL security group"
  vpc_id      = module.vpc.vpc_id

  # ingress
  ingress_with_cidr_blocks = [
    {
      from_port   = 5432
      to_port     = 5432
      protocol    = "tcp"
      description = "PostgreSQL access from within VPC"
      cidr_blocks = module.vpc.vpc_cidr_block
    },
  ]

  tags = merge(
    var.additional_tags,
    var.security_class_tags,
    {
      "dtit:sec:NetworkLayer" = "database"
    }
  )
}

module "lambda_security_group" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 4.0"

  name        = "${local.name}-lambda-security-group"
  description = "Complete Lambda security group"
  vpc_id      = module.vpc.vpc_id

  egress_with_cidr_blocks = [
    {
      from_port   = 0
      to_port     = 65535
      protocol    = "tcp"
      description = "Lambda access to the Internet from within the VPC"
      cidr_blocks = "0.0.0.0/0"
    },
  ]

  tags = merge(
    var.additional_tags,
    var.security_class_tags,
    {
      "dtit:sec:NetworkLayer" = "application"
    }
  )
}
