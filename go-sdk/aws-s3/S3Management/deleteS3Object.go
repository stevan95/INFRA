package S3Management

import (
	"aws-s3/ConfigS3Bucket"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

func DeleteObject(ctx context.Context, s3Client *s3.Client) error {
	objs, err := s3Client.ListObjects(ctx, &s3.ListObjectsInput{
		Bucket: aws.String(ConfigS3Bucket.S3Conf.BucketName),
		Prefix: aws.String(ConfigS3Bucket.S3Conf.KeyToDelete),
	})

	var objToDelete string

	for _, obj := range objs.Contents {
		objToDelete = string(*obj.Key)
		if err != nil {
			return fmt.Errorf("object not exists")
		} else {
			_, err := s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
				Bucket: aws.String(ConfigS3Bucket.S3Conf.BucketName),
				Key:    aws.String(objToDelete),
			})
			if err != nil {
				return fmt.Errorf("object not exists")
			}
		}
	}

	return nil
}
