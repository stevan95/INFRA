package awsinfo

import (
	awsconfig "aws-database-info/AWSConfig"
	"context"
	"database/sql"
	"fmt"
	"net/url"
)

func RetreiveRecords() {

	//Connect to the database
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

	ec2s := GetInstancesInfo(db)
	fmt.Println(len(ec2s))
	GetInstancesTags(db, &ec2s)

	for _, ec2 := range ec2s {
		fmt.Println("InstanceID: ", ec2.InstanceID)
		fmt.Println("PrivateIP: ", ec2.PrivateIPAddress)
		fmt.Println("PublicIP: ", ec2.PublicIpAddress)

		fmt.Println("EC2 Tags:")
		for key, value := range ec2.Tags {
			fmt.Println("  ", key, ":", value)
		}
	}
}

func GetInstancesInfo(db *sql.DB) []awsconfig.EC2Information {

	//Execute query on instances table to retrieve info
	rows, err := db.QueryContext(context.Background(), "SELECT InstanceID, PrivateIpAddress, PublicIpAddress FROM instances")
	if err != nil {
		fmt.Println("row.Scan", err)
		return nil
	}
	defer func() {
		_ = rows.Close()
	}()

	if rows.Err() != nil {
		fmt.Println("row.Error()", err)
		return nil
	}

	var slice []awsconfig.EC2Information
	for rows.Next() {
		instance := awsconfig.EC2Information{
			Tags: make(map[string]string),
		}
		if err := rows.Scan(&instance.InstanceID, &instance.PrivateIPAddress, &instance.PublicIpAddress); err != nil {
			fmt.Println("rows.Scan", err)
			return nil
		}

		slice = append(slice, instance)
	}

	return slice
}

func GetInstancesTags(db *sql.DB, slice *[]awsconfig.EC2Information) {
	// Retrieve tags for all instances
	rows, err := db.QueryContext(context.Background(), "SELECT InstanceID, TagsKey, TagsName FROM instances_tags")
	if err != nil {
		fmt.Println("row.Scan", err)
		return
	}
	defer func() {
		_ = rows.Close()
	}()

	if rows.Err() != nil {
		fmt.Println("row.Error()", err)
		return
	}

	// Map instance ID to tags
	tagsMap := make(map[string]map[string]string)
	for rows.Next() {
		var instanceID, key, value string
		if err := rows.Scan(&instanceID, &key, &value); err != nil {
			fmt.Println("rows.Scan", err)
			return
		}
		if _, ok := tagsMap[instanceID]; !ok {
			tagsMap[instanceID] = make(map[string]string)
		}
		tagsMap[instanceID][key] = value
	}

	// Update instances with tags
	for i, instance := range *slice {
		if tags, ok := tagsMap[instance.InstanceID]; ok {
			instance.Tags = tags
			(*slice)[i] = instance
		}
	}
}
