#!/bin/bash

set -e

# Function to run tests
run_tests() {
    echo "Running tests..."
    cd ./src
    go test ./cmd/desc-ec2/...
    cd ..
}

# Function to build the project
build_project() {
    echo "Building project..."
    cd ./src
    go build -o ../bin/desc-ec2/bootstrap ./cmd/desc-ec2
    cd ..

    zip -j bin/desc-ec.zip bin/desc-ec2/bootstrap
}

# Function to deploy infrastructure using Terraform
deploy_infrastructure() {
    echo "Deploying infrastructure..."
    cd infra
    terraform init
    terraform apply -auto-approve
    cd ..
}

# Main script execution
main() {
    run_tests
    build_project
    deploy_infrastructure
    echo "CI/CD pipeline completed successfully."
}

main