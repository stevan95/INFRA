package S3Management

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"aws-s3/ConfigS3Bucket"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

func UploadToS3Bucket(ctx context.Context, s3Client *s3.Client) error {
	//Open given location
	file, err := os.Open(ConfigS3Bucket.S3Conf.PathToUpload)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	defer file.Close()

	//Get file status
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	// If provided path is not directory
	if !fileInfo.IsDir() {
		file_to_upload, err := ioutil.ReadFile(ConfigS3Bucket.S3Conf.PathToUpload)
		if err != nil {
			fmt.Printf("Filed to read file: %s", err)
		}
		//Upload the file
		uploader := manager.NewUploader(s3Client)
		_, err = uploader.Upload(ctx, &s3.PutObjectInput{
			Bucket: aws.String(ConfigS3Bucket.S3Conf.BucketName),
			Key:    aws.String("steva/" + ConfigS3Bucket.S3Conf.PathToUpload),
			Body:   bytes.NewReader(file_to_upload),
		})
		if err != nil {
			fmt.Printf("Filed to uplad object error: %s", err)
		}
	} else {
		dir, err := os.ReadDir(ConfigS3Bucket.S3Conf.PathToUpload)
		if err != nil {
			log.Fatal(err)
		}

		var filenames []string
		for _, entry := range dir {
			filenames = append(filenames, ConfigS3Bucket.S3Conf.PathToUpload+"/"+entry.Name())
		}

		for _, file := range filenames {
			//Upload files from directory
			file_to_upload, err := ioutil.ReadFile(file)
			if err != nil {
				fmt.Printf("Filed to read file: %s", err)
			}
			uploader := manager.NewUploader(s3Client)
			_, err = uploader.Upload(ctx, &s3.PutObjectInput{
				Bucket: aws.String(ConfigS3Bucket.S3Conf.BucketName),
				Key:    aws.String(file),
				Body:   bytes.NewReader(file_to_upload),
			})
			fmt.Printf("Upload: %s\n", file_to_upload)

			if err != nil {
				fmt.Printf("Filed to uplad object error: %s", err)
			}
		}
	}

	return nil
}
