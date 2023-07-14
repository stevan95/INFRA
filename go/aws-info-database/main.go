package main

import (
	awsinfo "aws-database-info/AWSinfo"
	"context"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	ctx := context.Background()
	awsinfo.GetListofAllInstances(ctx, "us-east-1")
	awsinfo.RetreiveRecords()
}
