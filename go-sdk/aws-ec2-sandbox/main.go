package main

import (
	"aws-ec2-sandbox/EC2Management"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func main() {

	var (
		instanceID    string
		err           error
		region        string
		action        string
		instanceIDarg string
	)

	ctx := context.Background()

	flag.StringVar(&region, "region", "", "Define Region where you want to run your EC2 Instance.")
	flag.StringVar(&action, "action", "", "Define which action you want to perform Create/Stop/Terminate EC2 Instance.")
	flag.StringVar(&instanceIDarg, "instanceID", "", "Set instance ID of machine that you want to manage.")
	flag.Parse()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -region=<aws_region> -action=<create,start,stop,terminate> -terminateInstanceID=<instance_id> -instanceID=<instance_id>\n", os.Args[0])
		flag.PrintDefaults()
	}

	if region == "" {
		fmt.Print("Error, Region is mandatory argument not specified.\n")
		flag.Usage()
		os.Exit(1)
	}

	if action == "" {
		fmt.Print("Error, action is mandatory argument not specified.\n")
		flag.Usage()
		os.Exit(1)
	}

	switch action {
	case "create":
		if instanceID, err = EC2Management.CreateEC2(ctx, "us-east-1"); err != nil {
			fmt.Printf("CreateEC2 error: %s", err)
			os.Exit(1)
		}

		fmt.Printf("Instance ID: %s", instanceID)
	case "terminate":
		if instanceIDarg == "" {
			fmt.Print("Error, terminateInstanceID is mandatory argument not specified.\n")
			flag.Usage()
			os.Exit(1)
		}

		if instanceID, err = EC2Management.TerminateEC2(ctx, "us-east-1", instanceIDarg); err != nil {
			fmt.Printf("Terminated ec2 error: %s", err)
			os.Exit(1)
		}

		fmt.Printf("Terminated Instance ID: %s", instanceID)
	case "stop":
		if instanceIDarg == "" {
			fmt.Print("Error, instanceID is mandatory argument not specified.\n")
			flag.Usage()
			os.Exit(1)
		}

		if instanceID, err = EC2Management.StopEC2Instance(ctx, "us-east-1", instanceIDarg); err != nil {
			fmt.Printf("Stopped ec2 error: %s", err)
			os.Exit(1)
		}
		fmt.Printf("Stopped Instance ID: %s", instanceID)
	case "start":
		if instanceIDarg == "" {
			fmt.Print("Error, instanceID is mandatory argument not specified.\n")
			flag.Usage()
			os.Exit(1)
		}

		if instanceID, err = EC2Management.StartStoppedEC2Instance(ctx, "us-east-1", instanceIDarg); err != nil {
			fmt.Printf("Stopped ec2 error: %s", err)
			os.Exit(1)
		}

		fmt.Printf("Started Instance ID: %s", instanceID)

	default:
		fmt.Print("Wrong action is specified possible values are create/start/stop/terminate.")
	}

	GetListofAllInstances(ctx, "us-east-1")

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

	if len(EC2List.Reservations) == 0 {
		fmt.Printf("No running instances in selected region.")
	}

	type EC2Information struct {
		InstanceID       string            `json:"InstanceID"`
		PrivateIPAddress string            `json:"PrivateIpAddress"`
		PrivateDnsName   string            `json:"PrivateDnsName"`
		SubnetId         string            `json:"SubnetId"`
		PublicIpAddress  string            `json:"PublicIpAddress"`
		InstanceType     string            `json:"InstanceType"`
		Tags             map[string]string `json:"Tags"`
	}

	tagMap := make(map[string]string)
	for idx, _ := range EC2List.Reservations {
		for _, inst := range EC2List.Reservations[idx].Instances {

			for _, tag := range inst.Tags {
				tagMap[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
			}

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
				PrivateIPAddress: *inst.PrivateIpAddress,
				PrivateDnsName:   Private_DNS_Name,
				SubnetId:         *inst.SubnetId,
				PublicIpAddress:  Public_Ip_Address,
				InstanceType:     string(inst.InstanceType),
				Tags:             tagMap,
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
