data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

resource "aws_apigatewayv2_api" "this" {
  name          = var.api_gateway_name
  protocol_type = "HTTP"

  tags = var.api_gateway_tags
}

resource "aws_apigatewayv2_integration" "this" {
  api_id           = aws_apigatewayv2_api.this.id
  integration_type = "AWS_PROXY"

  connection_type    = "INTERNET"
  integration_method = var.api_gateway_http_method
  integration_uri    = var.api_gateway_integration_uri
}

resource "aws_apigatewayv2_stage" "this" {
  api_id = aws_apigatewayv2_api.this.id
  name   = var.api_gateway_stage_name

  auto_deploy = true

  access_log_settings {
    // TODO: Change to variable
    destination_arn = "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:DevSecOps_APIGateway_Access_Logging"
    //destination_arn = "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws/apigateway/${var.api_gateway_name}"
    format = jsonencode(
      {
        caller         = "$context.identity.caller"
        httpMethod     = "$context.httpMethod"
        ip             = "$context.identity.sourceIp"
        protocol       = "$context.protocol"
        requestId      = "$context.requestId"
        requestTime    = "$context.requestTime"
        resourcePath   = "$context.resourcePath"
        responseLength = "$context.responseLength"
        status         = "$context.status"
        user           = "$context.identity.user"
      }
    )
  }

  default_route_settings {
    data_trace_enabled       = false
    detailed_metrics_enabled = false
    throttling_burst_limit   = var.api_gateway_throttling_burst_limit
    throttling_rate_limit    = var.api_gateway_throttling_rate_limit
  }

  stage_variables = {}

  tags = var.api_gateway_tags
}
