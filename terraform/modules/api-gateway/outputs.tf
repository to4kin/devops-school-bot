output "api_gateway_arn" {
  description = "The ARN of the API gateway"
  value       = aws_apigatewayv2_api.this.arn
}

output "api_gateway_id" {
  description = "API identifier"
  value       = aws_apigatewayv2_api.this.id
}

output "api_gateway_stage_execution_arn" {
  description = "ARN prefix to be used in an aws_lambda_permission's source_arn attribute"
  value       = aws_apigatewayv2_stage.this.execution_arn
}

output "api_gateway_stage_invoke_url" {
  description = "URL to invoke the API pointing to the stage"
  value       = aws_apigatewayv2_stage.this.invoke_url
}

output "api_gateway_stage_id" {
  description = "Stage identifier"
  value       = aws_apigatewayv2_stage.this.id
}
