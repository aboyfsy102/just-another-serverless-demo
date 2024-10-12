variable "lambda_function_name" {
  description = "The name of the Lambda function"
  type        = string
  default     = "describe_ec2"
}

variable "region" {
  description = "The AWS region"
  type        = string
  default     = "ap-southeast-1"
}