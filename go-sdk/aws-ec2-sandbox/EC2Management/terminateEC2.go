package EC2Management

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func TerminateEC2(ctx context.Context, region string, instance_id string) (string, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		fmt.Printf("unable to load sdk config, %v", err)
	}

	ec2Client := ec2.NewFromConfig(cfg)

	instanceIDSlice := []string{instance_id}

	_, err = ec2Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: instanceIDSlice,
	})
	if err != nil {
		return "", fmt.Errorf("no such instance, wrong ID, %v", err)
	}

	_, err = ec2Client.TerminateInstances(ctx, &ec2.TerminateInstancesInput{
		InstanceIds: instanceIDSlice,
	})

	if err != nil {
		return "", fmt.Errorf("unable to terminate instance, %v", err)
	}

	return instance_id, nil
}
