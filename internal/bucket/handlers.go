package bucket

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

type BucketObject struct {
	Spec struct {
		Name   string `json:"name,omitempty"`
		Region string `json:"region,omitempty"`
	} `json:"spec,omitempty"`
}

func (c *Client) AddResourceHandler(obj any) {
	nq, ok := obj.(*unstructured.Unstructured)
	if !ok {
		log.Printf("error in formatting of object\n")
	}

	var bo BucketObject
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(nq.Object, &bo)
	if err != nil {
		log.Printf("error converting from unstructured to NQObject: %v\n", err)
	}

	_, err = c.api.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(bo.Spec.Name),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(bo.Spec.Region),
		},
	})
	if err != nil {
		log.Printf("error creating bucket: %v\n", err)
	}
}

func (c *Client) UpdateResourceHandler(oldObj, newObj any) {
	log.Printf("update hit\n")
}

func (c *Client) DeleteResourceHandler(obj any) {
	bucket := obj.(*unstructured.Unstructured)
	name, ok, err := unstructured.NestedString(bucket.Object, "spec", "name")
	if err != nil || !ok {
		log.Printf("error finding nested string in unstructed obj\n")
	}

	c.api.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(name),
	})

	_, err = c.api.DeleteBucket(context.TODO(), &s3.DeleteBucketInput{
		Bucket: aws.String(name),
	})
	if err != nil {
		log.Printf("error deleting bucket: %v\n", err)
	}
}
