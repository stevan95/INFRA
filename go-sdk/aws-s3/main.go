package main

import (
	"aws-s3/ConfigS3Bucket"
	"aws-s3/S3Management"
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {

	ConfigS3Bucket.S3Conf.InitS3Config("aws-s3-test-steva-95", "us-west-2", "test_files", "test_files/file1.txt")

	var (
		s3Client   *s3.Client
		err        error
		downloaded []byte
	)
	ctx := context.Background()

	if s3Client, err = initS3Client(ctx); err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	if err = S3Management.CreateS3Bucket(ctx, s3Client); err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
	if err = S3Management.UploadToS3Bucket(ctx, s3Client); err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
	fmt.Printf("File %s is uploaded.", ConfigS3Bucket.S3Conf.PathToUpload)

	if downloaded, err = S3Management.GetFileFromS3(ctx, s3Client); err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
	fmt.Printf("Downloaded file %s.", downloaded)
}

func initS3Client(ctx context.Context) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(ConfigS3Bucket.S3Conf.RegionName))
	if err != nil {
		return nil, fmt.Errorf("error: %s", err)
	}

	return s3.NewFromConfig(cfg), nil
}
