package main

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/stretchr/testify/assert"
)

type mockEC2Client struct {
	ec2iface.EC2API
}

func (m *mockEC2Client) DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return &ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{
			{
				Instances: []*ec2.Instance{
					{
						InstanceId:       aws.String("i-1234567890abcdef0"),
						State:            &ec2.InstanceState{Name: aws.String("running")},
						PrivateIpAddress: aws.String("192.168.1.1"),
						LaunchTime:       aws.Time(time.Now().UTC()),
						ImageId:          aws.String("ami-12345678"),
						Tags: []*ec2.Tag{
							{
								Key:   aws.String("Name"),
								Value: aws.String("TestInstance"),
							},
						},
					},
				},
			},
		},
	}, nil
}

func TestHandler(t *testing.T) {
	mockSvc := &mockEC2Client{}

	request := events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodGet,
	}

	expectedInstances := []InstanceInfo{
		{
			InstanceID:       "i-1234567890abcdef0",
			TagName:          "TestInstance",
			InstanceState:    "running",
			PrivateIPAddress: "192.168.1.1",
			LaunchTime:       time.Now().UTC().Format(time.RFC3339),
			ImageID:          "ami-12345678",
		},
	}

	expectedBody, _ := json.Marshal(expectedInstances)

	response, err := handler(context.Background(), request, mockSvc)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.JSONEq(t, string(expectedBody), response.Body)
}
