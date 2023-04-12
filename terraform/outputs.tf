output "db_instance_address" {
  description = "The address of the RDS instance"
  value       = module.db.db_instance_address
}

output "db_instance_availability_zone" {
  description = "The availability zone of the RDS instance"
  value       = module.db.db_instance_availability_zone
}

output "db_instance_endpoint" {
  description = "The connection endpoint"
  value       = module.db.db_instance_endpoint
}

output "db_instance_engine_version_actual" {
  description = "The running version of the database"
  value       = module.db.db_instance_engine_version_actual
}

output "db_instance_status" {
  description = "The RDS instance status"
  value       = module.db.db_instance_status
}

output "db_instance_name" {
  description = "The database name"
  value       = module.db.db_instance_name
}

output "db_instance_username" {
  description = "The master username for the database"
  value       = module.db.db_instance_username
  sensitive   = true
}

output "db_instance_password" {
  description = "The database password (this password may be old, because Terraform doesn't track it after initial creation)"
  value       = module.db.db_instance_password
  sensitive   = true
}

output "db_instance_port" {
  description = "The database port"
  value       = module.db.db_instance_port
}

# Docker Image
output "docker_image_uri" {
  description = "The ECR Docker image URI used to deploy Lambda Function"
  value       = module.docker_image.image_uri
}

# API Gateway
output "api_gateway_arn" {
  description = "The ARN of the API gateway"
  value       = module.api_gateway.api_gateway_arn
}

output "api_gateway_id" {
  description = "Stage identifier"
  value       = module.api_gateway.api_gateway_id
}

output "api_gateway_stage_execution_arn" {
  description = "ARN prefix to be used in an aws_lambda_permission's source_arn attribute"
  value       = module.api_gateway.api_gateway_stage_execution_arn
}

output "api_gateway_stage_invoke_url" {
  description = "URL to invoke the API pointing to the stage"
  value       = module.api_gateway.api_gateway_stage_invoke_url
}

output "api_gateway_stage_id" {
  description = "API identifier"
  value       = module.api_gateway.api_gateway_stage_id
}

# Lambda Function
# output "lambda_function_arn" {
#   description = "The ARN of the Lambda Function"
#   value       = module.lambda_function.lambda_function_arn
# }

