package S3Management

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"aws-s3/ConfigS3Bucket"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

func GetFileFromS3(ctx context.Context, s3Client *s3.Client) ([]byte, error) {
	downloader := manager.NewDownloader(s3Client)

	buffer := manager.NewWriteAtBuffer([]byte{})

	numBytes, err := downloader.Download(ctx, buffer, &s3.GetObjectInput{
		Bucket: aws.String(ConfigS3Bucket.S3Conf.BucketName),
		Key:    aws.String(ConfigS3Bucket.S3Conf.FileToDownload),
	})
	if err != nil {
		fmt.Printf("Filed to downlaod object: %s", err)
	}

	if numBytesRecived := len(buffer.Bytes()); numBytes != int64(numBytesRecived) {
		fmt.Printf("Filed number of bytes recived doesnt match: %d vs %d", numBytes, numBytesRecived)
	}

	filetoDownloadPaths := strings.Split(ConfigS3Bucket.S3Conf.FileToDownload, "/")
	filetoDownloadName := filetoDownloadPaths[len(filetoDownloadPaths)-1]
	err = ioutil.WriteFile(filetoDownloadName, buffer.Bytes(), 0644)
	if err != nil {
		fmt.Println("failed to write file:", err)
		return []byte{}, nil
	}

	return buffer.Bytes(), nil
}
