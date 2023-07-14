package main

import (
	"aws-s3/ConfigS3Bucket"
	"aws-s3/S3Management"
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {

	var (
		s3BucketName      string
		regionName        string
		fileToUpload      string
		fileToDownload    string
		keyObjectToDelete string
		action            string
	)

	flag.StringVar(&s3BucketName, "bucketName", "", "Define name of your S3 bucket.")
	flag.StringVar(&action, "action", "", "Define which action you want to perform createS3Bucket/getFile/uploadFile/deleteObject.")
	flag.StringVar(&regionName, "region", "", "Set region name.")
	flag.StringVar(&fileToUpload, "fileToUpload", "", "Set name of the file which you want to upload.")
	flag.StringVar(&fileToDownload, "fileToDownload", "", "Set the neame of the file which you want to download.")
	flag.StringVar(&keyObjectToDelete, "keyObjectToDelete", "", "Set the key of object which you want to delete.")
	flag.Parse()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			"Usage: %s -region=<aws_region> -bucketName=<bucket_name> -action=<createS3Bucket/deleteS3Bucket/getFile/uploadFile/deleteObject>\nOptional Args depending on action:\n-fileToUpload=<fileToUpload>\n-fileToDownload=<fileToDownload>\n-keyObjectToDelete<keyObjectToDelete>", os.Args[0])
		flag.PrintDefaults()
	}

	if s3BucketName == "" {
		fmt.Print("Error, Bucket name is mandatory argument not specified.\n")
		flag.Usage()
		os.Exit(1)
	}

	if regionName == "" {
		fmt.Print("Error, Bucket name is mandatory argument not specified.\n")
		flag.Usage()
		os.Exit(1)
	}

	if action == "" {
		fmt.Print("Error, action is mandatory argument, you need to specify action.\n")
		flag.Usage()
		os.Exit(1)
	}

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

	switch action {
	case "createS3Bucket":
		ConfigS3Bucket.S3Conf.InitS3Config(s3BucketName, regionName, nil, nil, nil)
		if err = S3Management.CreateS3Bucket(ctx, s3Client); err != nil {
			fmt.Printf("Error: %s", err)
			os.Exit(1)
		}

		if err == nil {
			fmt.Printf("Bucket Name: %s", ConfigS3Bucket.S3Conf.BucketName)
		}

	case "uploadFile":
		ConfigS3Bucket.S3Conf.InitS3Config(s3BucketName, regionName, &fileToUpload, nil, nil)

		if err = S3Management.UploadToS3Bucket(ctx, s3Client); err != nil {
			fmt.Printf("Error: %s", err)
			os.Exit(1)
		}
		fmt.Printf("File %s is uploaded.", ConfigS3Bucket.S3Conf.PathToUpload)

	case "getFile":
		ConfigS3Bucket.S3Conf.InitS3Config(s3BucketName, regionName, nil, &fileToDownload, nil)

		if downloaded, err = S3Management.GetFileFromS3(ctx, s3Client); err != nil {
			fmt.Printf("Error: %s", err)
			os.Exit(1)
		}
		fmt.Printf("Downloaded file %s.", downloaded)

	case "deleteS3Bucket":
		ConfigS3Bucket.S3Conf.InitS3Config(s3BucketName, regionName, nil, nil, nil)

		if err = S3Management.DeleteS3Bucket(ctx, s3Client); err != nil {
			fmt.Printf("Error: %s", err)
			os.Exit(1)
		}
		fmt.Printf("Deleted Bucket: %s.", ConfigS3Bucket.S3Conf.BucketName)

	case "deleteObject":
		ConfigS3Bucket.S3Conf.InitS3Config(s3BucketName, regionName, nil, nil, &keyObjectToDelete)

		if err = S3Management.DeleteObject(ctx, s3Client); err != nil {
			fmt.Printf("Error: %s", err)
			os.Exit(1)
		}
		fmt.Printf("Deleted Object: %s.", ConfigS3Bucket.S3Conf.KeyToDelete)

	default:
		fmt.Print("Wrong action is specified possible values are createS3Bucket/deleteS3Bucket/getFile/uploadFile/deleteObject.")
	}

}

func initS3Client(ctx context.Context) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(ConfigS3Bucket.S3Conf.RegionName))
	if err != nil {
		return nil, fmt.Errorf("error: %s", err)
	}

	return s3.NewFromConfig(cfg), nil
}
