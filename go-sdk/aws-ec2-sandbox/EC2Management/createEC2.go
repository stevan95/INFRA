package EC2Management

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go/aws"
)

func CreateEC2(ctx context.Context, region string) (string, error) {

	var (
		securityGroup *ec2.CreateSecurityGroupOutput
		sgIDs         []string
	)

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return "", fmt.Errorf("unable to load sdk config, %v", err)
	}

	ec2Client := ec2.NewFromConfig(cfg)

	keyPairs, err := ec2Client.DescribeKeyPairs(ctx, &ec2.DescribeKeyPairsInput{
		KeyNames: []string{"go-aws-demo"},
	})

	if keyPairs == nil || len(keyPairs.KeyPairs) == 0 {
		keyPairs, err := ec2Client.CreateKeyPair(ctx, &ec2.CreateKeyPairInput{
			KeyName: aws.String("go-aws-demo"),
		})

		if err != nil && !strings.Contains(err.Error(), "InvalidKeyPair.NotFound") {
			return "", fmt.Errorf("CreateKeyPair error, %s", err)
		}

		err = os.WriteFile("go-aws-ec2.pem", []byte(*keyPairs.KeyMaterial), 0400)
		if err != nil {
			return "", fmt.Errorf("WriteFile error, %s", err)
		}
	}
	if err != nil {
		return "", fmt.Errorf("DescribeImages error, %s", err)
	}

	vpc_ids, err := ec2Client.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{
		Filters: []types.Filter{
			{
				Name: aws.String("isDefault"),
				Values: []string{
					"true",
				},
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("DescribeVpcs error, %s", err)
	}

	defaultVPC := vpc_ids.Vpcs[0]

	describe_sg, _ := ec2Client.DescribeSecurityGroups(ctx, &ec2.DescribeSecurityGroupsInput{
		GroupNames: []string{"Test-EC2"},
	})

	if describe_sg == nil || len(describe_sg.SecurityGroups) == 0 {
		securityGroup, err = ec2Client.CreateSecurityGroup(ctx, &ec2.CreateSecurityGroupInput{
			GroupName:   aws.String("Test-EC2"),
			Description: aws.String("EC2-Group which allow SSH on port 22."),
			VpcId:       defaultVPC.VpcId,
		})
		if err != nil {
			return "", fmt.Errorf("CreateSecurityGroup error, %s", err)
		}

		_, err = ec2Client.AuthorizeSecurityGroupIngress(ctx, &ec2.AuthorizeSecurityGroupIngressInput{
			GroupId: securityGroup.GroupId,
			IpPermissions: []types.IpPermission{
				{
					IpProtocol: aws.String("tcp"),
					FromPort:   aws.Int32(22),
					ToPort:     aws.Int32(22),
					IpRanges: []types.IpRange{
						{
							CidrIp: aws.String("0.0.0.0/0"),
						},
					},
				},
			},
		})
		if err != nil {
			return "", fmt.Errorf("AuthorizeSecurityGroupIngress, %s", err)
		}

		sgIDs = []string{*securityGroup.GroupId} //Dereference pointer to get slice of strings.
	} else {
		sgIDs = []string{*describe_sg.SecurityGroups[0].GroupId}
	}

	imageOutput, err := ec2Client.DescribeImages(ctx, &ec2.DescribeImagesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("name"),
				Values: []string{"ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"},
			},
			{
				Name:   aws.String("virtualization-type"),
				Values: []string{"hvm"},
			},
		},
		Owners: []string{"099720109477"},
	})
	if err != nil {
		return "", fmt.Errorf("DescribeImages error, %s", err)
	}

	if len(imageOutput.Images) == 0 {
		return "", fmt.Errorf("imageOutput.Images is empty, %s", err)
	}

	instance, err := ec2Client.RunInstances(ctx, &ec2.RunInstancesInput{
		ImageId:          imageOutput.Images[0].ImageId,
		SecurityGroupIds: sgIDs,
		KeyName:          aws.String("go-aws-demo"),
		InstanceType:     types.InstanceTypeT3Micro,
		MinCount:         aws.Int32(1),
		MaxCount:         aws.Int32(1),
	})

	if err != nil {
		return "", fmt.Errorf("run instance error, %s", err)
	}

	if len(instance.Instances) == 0 {
		return "", fmt.Errorf("instance.Instances is zero lenght, %s", err)
	}

	return *instance.Instances[0].InstanceId, nil
}
