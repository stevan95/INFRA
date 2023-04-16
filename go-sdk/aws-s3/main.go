package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const bucketName = "aws-s3-test-steva-95"
const regionName = "us-west-2"
const pathToUpload = "test_files"
const filetoDownload = "test_files/file1.txt"

func main() {
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
	if err = createS3Bucket(ctx, s3Client); err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
	/*if err = uploadToS3Bucket(ctx, s3Client); err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
	fmt.Printf("File %s is uploaded.", pathToUpload)*/

	if downloaded, err = getFileFromS3(ctx, s3Client); err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
	fmt.Printf("Downloaded file %s.", downloaded)
}

func initS3Client(ctx context.Context) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(regionName))
	if err != nil {
		return nil, fmt.Errorf("error: %s", err)
	}

	return s3.NewFromConfig(cfg), nil
}

func createS3Bucket(ctx context.Context, s3Client *s3.Client) error {
	allBuckets, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		fmt.Printf("Listing buckets error: %s", err)
	}
	found := false

	for _, bucket := range allBuckets.Buckets {
		if *bucket.Name == bucketName {
			found = true
		}
	}

	if !found {
		_, err := s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
			CreateBucketConfiguration: &types.CreateBucketConfiguration{
				LocationConstraint: regionName,
			},
		})
		if err != nil {
			fmt.Printf("Creating bucket error: %s", err)
		}
	}

	return nil
}

func uploadToS3Bucket(ctx context.Context, s3Client *s3.Client) error {
	//Open given location
	file, err := os.Open(pathToUpload)
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
		file_to_upload, err := ioutil.ReadFile(pathToUpload)
		if err != nil {
			fmt.Printf("Filed to read file: %s", err)
		}
		//Upload the file
		uploader := manager.NewUploader(s3Client)
		_, err = uploader.Upload(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("steva/" + pathToUpload),
			Body:   bytes.NewReader(file_to_upload),
		})
		if err != nil {
			fmt.Printf("Filed to uplad object error: %s", err)
		}
	} else {
		dir, err := os.ReadDir(pathToUpload)
		if err != nil {
			log.Fatal(err)
		}

		var filenames []string
		for _, entry := range dir {
			filenames = append(filenames, pathToUpload+"/"+entry.Name())
		}

		for _, file := range filenames {
			//Upload files from directory
			file_to_upload, err := ioutil.ReadFile(file)
			if err != nil {
				fmt.Printf("Filed to read file: %s", err)
			}
			uploader := manager.NewUploader(s3Client)
			_, err = uploader.Upload(ctx, &s3.PutObjectInput{
				Bucket: aws.String(bucketName),
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

func getFileFromS3(ctx context.Context, s3Client *s3.Client) ([]byte, error) {
	downloader := manager.NewDownloader(s3Client)

	buffer := manager.NewWriteAtBuffer([]byte{})

	numBytes, err := downloader.Download(ctx, buffer, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filetoDownload),
	})
	if err != nil {
		fmt.Printf("Filed to downlaod object: %s", err)
	}

	if numBytesRecived := len(buffer.Bytes()); numBytes != int64(numBytesRecived) {
		fmt.Printf("Filed number of bytes recived doesnt match: %d vs %d", numBytes, numBytesRecived)
	}

	filetoDownloadPaths := strings.Split(filetoDownload, "/")
	filetoDownloadName := filetoDownloadPaths[len(filetoDownloadPaths)-1]
	err = ioutil.WriteFile(filetoDownloadName, buffer.Bytes(), 0644)
	if err != nil {
		fmt.Println("failed to write file:", err)
		return []byte{}, nil
	}

	return buffer.Bytes(), nil
}
