package EC2Management

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func StopEC2Instance(ctx context.Context, region string, instanceID string) (string, error) {
	//Variables
	var (
		isRunning bool
	)

	//Create go sdk client to connect the API
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return "", fmt.Errorf("unable to load sdk config, %v", err)
	}

	instanceIDSlice := []string{instanceID}

	ec2Client := ec2.NewFromConfig(cfg)

	//Check if instance exist (if instance is not found throw error)
	_, err = ec2Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: instanceIDSlice,
	})
	if err != nil {
		return "", fmt.Errorf("no instance found error, %v", err)
	}

	//Get state of the instance
	ec2_status, err := ec2Client.DescribeInstanceStatus(ctx, &ec2.DescribeInstanceStatusInput{
		InstanceIds: instanceIDSlice,
	})
	if err != nil {
		return "", fmt.Errorf("describeInstanceStatus, %v", err)
	}

	for _, instanceStatus := range ec2_status.InstanceStatuses {
		if instanceStatus.InstanceState.Name == "Stopped" {
			isRunning = false
		} else {
			isRunning = true
		}
	}

	if isRunning {
		_, err = ec2Client.StopInstances(ctx, &ec2.StopInstancesInput{
			InstanceIds: instanceIDSlice,
		})
		if err != nil {
			return "", fmt.Errorf("StopInstances error, %s", err)
		}

		return instanceID, nil
	}

	return "Instance is already stopped.", nil
}
