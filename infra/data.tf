# Query the default VPC in the target region
data "aws_vpc" "default" {
  default = true
}

# Query the subnets in the default VPC
data "aws_subnets" "default_vpc_subnets" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }
}

# Output the list of subnet IDs
output "default_vpc_subnet_ids" {
  description = "List of subnet IDs in the default VPC"
  value       = data.aws_subnets.default_vpc_subnets.ids
}