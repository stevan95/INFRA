package SGManagement

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func RemoveSG(ctx context.Context, sg_id string) (string, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-west-2"))
	if err != nil {
		fmt.Printf("unable to load sdk config, %v", err)
	}

	awsClient := ec2.NewFromConfig(cfg)

	describe_sg, _ := awsClient.DescribeSecurityGroups(ctx, &ec2.DescribeSecurityGroupsInput{
		GroupIds: []string{sg_id},
	})

	if describe_sg == nil || len(describe_sg.SecurityGroups) == 0 {
		return "SG not exists", nil
	} else {
		awsClient.DeleteSecurityGroup(ctx, &ec2.DeleteSecurityGroupInput{
			GroupId: aws.String(sg_id),
		})

		return sg_id, nil
	}
}
