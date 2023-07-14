package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"aws-security-group/ConfigSG"
	"aws-security-group/SGManagement"
)

func main() {

	var (
		filename string
		action   string
		sgID     string
	)

	flag.StringVar(&filename, "filename", "", "File name which contain sg definitions.")
	flag.StringVar(&action, "action", "", "Define which action you want to perform create/delete")
	flag.StringVar(&sgID, "sgID", "", "Name of sg which you want to delete.")
	flag.Parse()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			"Usage: %s -action=<create/delete> -filename=<filename.json> -sgID<sgID(only when perform delete)>", os.Args[0])
		flag.PrintDefaults()
	}

	if action == "" {
		fmt.Print("action is required filed.\n")
		flag.Usage()
		os.Exit(1)
	}

	if filename == "" && action == "create" {
		fmt.Print("error, file name is not specified.\n")
		flag.Usage()
		os.Exit(1)
	} else if len(filename) != 0 && action != "create" {
		fmt.Print("filename requires create action.\n")
		flag.Usage()
		os.Exit(1)
	} else if len(filename) != 0 && action == "delete" {
		fmt.Print("delete action need other parameter (sgID).\n")
		flag.Usage()
		os.Exit(1)
	} else if len(sgID) != 0 && action != "delete" {
		fmt.Print("sgID reguires delete action.\n")
		flag.Usage()
		os.Exit(1)
	} else if sgID == "" && action == "delete" {
		fmt.Print("delete action requires sg.\n")
		flag.Usage()
		os.Exit(1)
	}

	switch action {
	case "create":
		file, err := ioutil.ReadFile(filename)
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
	case "delete":
		ctx := context.Background()

		rsgID, _ := SGManagement.RemoveSG(ctx, sgID)
		fmt.Printf("Removed Security Group: %s\n", rsgID)
	}

}
