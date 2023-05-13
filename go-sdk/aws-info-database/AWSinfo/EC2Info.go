package awsinfo

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

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
		fmt.Printf("DescribeInstances error, %s", err)
	}

	if len(EC2List.Reservations) == 0 {
		fmt.Printf("No running instances in selected region.")
		return
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

			//INSERT Record INTO Database if it not exist
			dsn := url.URL{
				Scheme: "postgres",
				Host:   "192.168.8.105:5432",
				User:   url.UserPassword("postgres", "mysecretpassword"),
				Path:   "ec2informations",
			}

			q := dsn.Query()
			q.Add("sslmode", "disable")

			dsn.RawQuery = q.Encode()

			db, err := sql.Open("pgx", dsn.String())
			if err != nil {
				fmt.Println("sql.Open", err)
				return
			}
			defer func() {
				_ = db.Close()
			}()

			var exists bool
			err = db.QueryRowContext(context.Background(),
				"SELECT EXISTS(SELECT InstanceID FROM instances WHERE InstanceID = $1)", inst.InstanceId).Scan(&exists)
			if err != nil {
				panic(err)
			}

			if !exists {
				_, err = db.ExecContext(context.Background(),
					"INSERT INTO instances(InstanceID, PrivateIpAddress, PrivateDnsName, SubnetId, PublicIpAddress, InstanceType) VALUES($1, $2, $3, $4, $5, $6)", inst.InstanceId, inst.PrivateIpAddress, Private_DNS_Name, inst.SubnetId, Public_Ip_Address, inst.InstanceType)

				if err != nil {
					fmt.Println("db.ExecContext: ", err)
					return
				}
			}

			//Insert tags into table
			for key, value := range tagMap {
				var existsTag bool
				err = db.QueryRowContext(context.Background(),
					"SELECT EXISTS(SELECT TagsName FROM instances_tags WHERE TagsName = $1)", value).Scan(&existsTag)
				if err != nil {
					panic(err)
				}

				if !existsTag {
					_, err = db.ExecContext(context.Background(),
						"INSERT INTO instances_tags(InstanceID, TagsKey, TagsName) VALUES($1, $2, $3)", inst.InstanceId, key, value)

					if err != nil {
						fmt.Println("db.ExecContext: ", err)
						return
					}
				}
			}
		}
	}
}
