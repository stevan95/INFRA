package main

import (
	"context"
	"fmt"

	"aws-security-group/SGManagement"
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
	sgInfo["sgPortsIPAddressProtocol"] = "0.0.0.0/0:22:tcp,109.245.79.198/32:8080:tcp"

	vpcID, _ := SGManagement.CreateSecurityGroup(ctx, sgInfo)

	fmt.Printf("VpcID: %s", vpcID)

}
