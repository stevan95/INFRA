package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"aws-security-group/ConfigSG"
	"aws-security-group/SGManagement"
)

func main() {

	file, err := ioutil.ReadFile("configFile.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var sgs map[string]ConfigSG.SGDetails
	err = json.Unmarshal(file, &sgs)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	for sgName, sgDetails := range sgs {
		ctx := context.Background()

		sgID, _ := SGManagement.CreateSecurityGroup(ctx, sgDetails)
		fmt.Printf("Security Group(%s) ID: %s\n", sgName, sgID)
	}

	/*rsgID, _ := SGManagement.RemoveSG(ctx, sgID)
	fmt.Printf("Removed Security Group: %s\n", rsgID)*/
}
