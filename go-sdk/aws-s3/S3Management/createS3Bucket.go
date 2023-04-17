package S3Management

import (
	"context"
	"fmt"

	"aws-s3/ConfigS3Bucket"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws"
)

func CreateS3Bucket(ctx context.Context, s3Client *s3.Client) error {
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

	if !found {
		_, err := s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(ConfigS3Bucket.S3Conf.BucketName),
			CreateBucketConfiguration: &types.CreateBucketConfiguration{
				LocationConstraint: types.BucketLocationConstraint(ConfigS3Bucket.S3Conf.RegionName),
			},
		})
		if err != nil {
			fmt.Printf("Creating bucket error: %s", err)
		}
	}

	return nil
}
