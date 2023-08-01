package bucket

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3API interface {
	CreateBucket(context.Context, *s3.CreateBucketInput, ...func(*s3.Options)) (*s3.CreateBucketOutput, error)
	DeleteBucket(context.Context, *s3.DeleteBucketInput, ...func(*s3.Options)) (*s3.DeleteBucketOutput, error)
}

func CreateBucket(ctx context.Context, client S3API, input *s3.CreateBucketInput) (*s3.CreateBucketOutput, error) {
	return client.CreateBucket(ctx, input)
}
func DeleteBucket(ctx context.Context, client S3API, input *s3.DeleteBucketInput) (*s3.DeleteBucketOutput, error) {
	return client.DeleteBucket(ctx, input)
}
