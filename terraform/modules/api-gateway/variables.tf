
variable "api_gateway_name" {
  description = "The name of the resource"
  type        = string
  default     = null
}

variable "api_gateway_http_method" {
  description = "HTTP Method (GET, POST, PUT, DELETE, HEAD, OPTIONS, ANY)"
  type        = string
  default     = "GET"
}

variable "api_gateway_integration_uri" {
  description = "The URI should be of the form arn:aws:apigateway:{region}:{subdomain.service|service}:{path|action}/{service_api}. region, subdomain and service are used to determine the right endpoint"
  type        = string
  default     = null
}

variable "api_gateway_stage_name" {
  description = "Name of the stage"
  type        = string
  default     = null
}

variable "api_gateway_throttling_burst_limit" {
  description = "Throttling burst limit. Default: 5000"
  type        = number
  default     = 5000
}

variable "api_gateway_throttling_rate_limit" {
  description = "Throttling rate limit. Default: 10000"
  type        = number
  default     = 10000
}

variable "api_gateway_tags" {
  description = "A map of tags to assign to API Gateway"
  type        = map(string)
  default     = {}
}
