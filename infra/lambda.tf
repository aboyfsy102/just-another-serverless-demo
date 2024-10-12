resource "aws_lambda_function" "describe_ec2" {
  function_name = var.lambda_function_name
  role          = aws_iam_role.lambda_exec_role.arn
  handler       = "index.handler"
  runtime       = "provided.al2"
  filename      = "../release/lambda_function_payload.zip"
#   source_code_hash = filebase64sha256("lambda_function_payload.zip")

  environment {
    variables = {
      AWS_REGION = var.region
    }
  }
}