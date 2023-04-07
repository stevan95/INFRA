package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// Create EC2 Instance +
// Stop EC2 instance +
// Terminate Instance +
// Start stopped Instance
// GetListofAllInstances +
func main() {

	var (
		instanceID string
		err        error
	)

	ctx := context.Background()

	/*if instanceID, err = createEC2(ctx, "us-east-1"); err != nil {
		fmt.Printf("CreateEC2 error: %s", err)
		os.Exit(1)
	}

	fmt.Printf("Instance ID: %s", instanceID)*/

	if instanceID, err = terminateEC2(ctx, "us-east-1", "i-081be3bc3c544e609"); err != nil {
		fmt.Printf("Terminated ec2 error: %s", err)
		os.Exit(1)
	}

	fmt.Printf("Terminated Instance ID: %s", instanceID)

	//GetListofAllInstances(ctx, "us-east-1")

	/*if instanceID, err = stopEC2Instance(ctx, "us-east-1", "i-081be3bc3c544e609"); err != nil {
		fmt.Printf("Stopped ec2 error: %s", err)
		os.Exit(1)
	}*/

	//fmt.Printf("Stopped Instance ID: %s", instanceID)

	/*if instanceID, err = startStoppedEC2Instance(ctx, "us-east-1", "i-081be3bc3c544e609"); err != nil {
		fmt.Printf("Stopped ec2 error: %s", err)
		os.Exit(1)
	}*/

	//fmt.Printf("Started Instance ID: %s", instanceID)

}

func createEC2(ctx context.Context, region string) (string, error) {

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

	describe_sg, err := ec2Client.DescribeSecurityGroups(ctx, &ec2.DescribeSecurityGroupsInput{
		GroupNames: []string{"Test-EC2"},
	})
	if err != nil {
		return "", fmt.Errorf("DescribeVpcs error, %s", err)
	}

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

func stopEC2Instance(ctx context.Context, region string, instanceID string) (string, error) {
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

func GetListofAllInstances(ctx context.Context, region string) {

	//Create go sdk client to connect the API
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		fmt.Printf("unable to load sdk config, %v", err)
	}

	ec2Client := ec2.NewFromConfig(cfg)

	EC2List, err := ec2Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name: aws.String("instance-state-name"),
				Values: []string{
					"running",
					"stopped",
				},
			},
		},
	})
	if err != nil {
		fmt.Printf("StopInstances error, %s", err)
	}

	type EC2Information struct {
		InstanceID       string `json:"InstanceID"`
		PrivateIPAddress string `json:"PrivateIpAddress"`
		PrivateDnsName   string `json:"PrivateDnsName"`
		SubnetId         string `json:"SubnetId"`
		PublicIpAddress  string `json:"PublicIpAddress"`
		InstanceType     string `json:"InstanceType"`
	}

	for idx, _ := range EC2List.Reservations {
		for _, inst := range EC2List.Reservations[idx].Instances {

			var (
				Private_DNS_Name  string
				Public_Ip_Address string
			)

			if inst.PrivateDnsName != nil {
				Private_DNS_Name = *inst.PrivateDnsName
			} else {
				Private_DNS_Name = "None"
			}

			if inst.PublicIpAddress != nil {
				Public_Ip_Address = *inst.PublicIpAddress
			} else {
				Public_Ip_Address = "None"
			}

			ec2_json, err := json.MarshalIndent(EC2Information{
				InstanceID:       *inst.InstanceId,
				PrivateIPAddress: *inst.PrivateDnsName,
				PrivateDnsName:   Private_DNS_Name,
				SubnetId:         *inst.SubnetId,
				PublicIpAddress:  Public_Ip_Address,
				InstanceType:     string(inst.InstanceType),
			}, "", "    ")
			if err != nil {
				fmt.Printf("unable to load sdk config, %v", err)
			}

			filename := *inst.InstanceId + "_" + "output.json"

			err = ioutil.WriteFile(filename, ec2_json, 0644)
			if err != nil {
				fmt.Printf("unable to print output as json, %v", err)
			}

		}
	}
}

func terminateEC2(ctx context.Context, region string, instance_id string) (string, error) {
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

func startStoppedEC2Instance(ctx context.Context, region string, instanceID string) (string, error) {
	//Variables
	var (
		isStopped bool
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
			isStopped = true
		} else {
			isStopped = false
		}
	}

	if !isStopped {
		_, err = ec2Client.StartInstances(ctx, &ec2.StartInstancesInput{
			InstanceIds: instanceIDSlice,
		})
		if err != nil {
			return "", fmt.Errorf("StartInstances error, %s", err)
		}

		return instanceID, nil
	}

	return "Instance is already Started.", nil
}
