package S3Management

import (
	"aws-s3/ConfigS3Bucket"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

func DeleteS3Bucket(ctx context.Context, s3Client *s3.Client) error {
	allBuckets, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		fmt.Printf("Listing buckets error: %s", err)
	}
	found := false

	for _, bucket := range allBuckets.Buckets {
		if *bucket.Name == ConfigS3Bucket.S3Conf.BucketName {
			found = true
		}
	}

	if found {
		_, err := s3Client.DeleteBucket(ctx, &s3.DeleteBucketInput{
			Bucket: aws.String(ConfigS3Bucket.S3Conf.BucketName),
		})
		if err != nil {
			fmt.Printf("Deleting bucket error: %s", err)
		}
	} else {
		return fmt.Errorf("bucket not exists")
	}

	return nil
}
