package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func main() {

	var (
		sgInfo = make(map[string]string)
	)

	ctx := context.Background()

	sgInfo["vpcName"] = "Default"
	sgInfo["sgName"] = "test_sg"
	sgInfo["sgDescription"] = "Allow SSH on port 22"
	sgInfo["sgProtocol"] = "tcp"
	sgInfo["sgPortsIPAddress"] = "0.0.0.0/0:22,109.245.79.198/32:8080"

	vpcID, _ := CreateSecurityGroup(ctx, sgInfo)

	fmt.Printf("VpcID: %s", vpcID)

}

func CreateSecurityGroup(ctx context.Context, sgInfo map[string]string) (string, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-west-2"))
	if err != nil {
		fmt.Printf("unable to load sdk config, %v", err)
	}

	awsClient := ec2.NewFromConfig(cfg)

	describe_vpc, _ := awsClient.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []string{sgInfo["vpcName"]},
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("DescribeVpcs error, %s", err)
	}

	VPCID := describe_vpc.Vpcs[0].VpcId

	describe_sg, _ := awsClient.DescribeSecurityGroups(ctx, &ec2.DescribeSecurityGroupsInput{
		GroupNames: []string{sgInfo["sgName"]},
	})

	var securityGroup *ec2.CreateSecurityGroupOutput
	if describe_sg == nil || len(describe_sg.SecurityGroups) == 0 {
		securityGroup, err = awsClient.CreateSecurityGroup(ctx, &ec2.CreateSecurityGroupInput{
			GroupName:   aws.String(sgInfo["sgName"]),
			Description: aws.String(sgInfo["sgDescription"]),
			VpcId:       VPCID,
		})
		if err != nil {
			return "", fmt.Errorf("CreateSecurityGroup error, %s", err)
		}

		address_ports := strings.Split(sgInfo["sgPortsIPAddress"], ",")
		for _, addr_port := range address_ports {

			addr := strings.Split(addr_port, ":")[0]
			port, _ := strconv.Atoi(strings.Split(addr_port, ":")[1])
			port32 := int32(port)

			_, err = awsClient.AuthorizeSecurityGroupIngress(ctx, &ec2.AuthorizeSecurityGroupIngressInput{
				GroupId: securityGroup.GroupId,
				IpPermissions: []types.IpPermission{
					{
						IpProtocol: aws.String(sgInfo["sgProtocol"]),
						FromPort:   aws.Int32(port32),
						ToPort:     aws.Int32(port32),
						IpRanges: []types.IpRange{
							{
								CidrIp: aws.String(addr),
							},
						},
					},
				},
			})
		}
		if err != nil {
			return "", fmt.Errorf("AuthorizeSecurityGroupIngress, %s", err)
		}

		return []string{*securityGroup.GroupId}[0], nil

	}

	return []string{*describe_sg.SecurityGroups[0].GroupId}[0], nil
}
