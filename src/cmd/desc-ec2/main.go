package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

type InstanceInfo struct {
	InstanceID       string `json:"instance_id"`
	TagName          string `json:"tag_name"`
	InstanceState    string `json:"instance_state"`
	PrivateIPAddress string `json:"private_ip_address"`
	LaunchTime       string `json:"launch_time"`
	ImageID          string `json:"image_id"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest, svc ec2iface.EC2API) (events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod != http.MethodGet {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       "Method not allowed",
		}, nil
	}

	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running")},
			},
		},
	}

	result, err := svc.DescribeInstances(input)
	if err != nil {
		log.Fatalf("Failed to describe instances: %v", err)
	}

	var instances []InstanceInfo
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			var tagName string
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					tagName = *tag.Value
					break
				}
			}

			instances = append(instances, InstanceInfo{
				InstanceID:       *instance.InstanceId,
				TagName:          tagName,
				InstanceState:    *instance.State.Name,
				PrivateIPAddress: *instance.PrivateIpAddress,
				LaunchTime:       instance.LaunchTime.UTC().Format(time.RFC3339),
				ImageID:          *instance.ImageId,
			})
		}
	}

	responseBody, err := json.Marshal(instances)
	if err != nil {
		log.Fatalf("Failed to marshal response: %v", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseBody),
	}, nil
}

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1"),
	})
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}

	svc := ec2.New(sess)
	lambda.Start(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return handler(ctx, request, svc)
	})
}
