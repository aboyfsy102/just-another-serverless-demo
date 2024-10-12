output "lambda_function_arn" {
  description = "The ARN of the Lambda function"
  value       = aws_lambda_function.describe_ec2.function_name
}

output "alb_dns_name" {
  description = "The DNS name of the ALB"
  value       = aws_lb.app_lb.dns_name
}