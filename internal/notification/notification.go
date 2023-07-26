package notification

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// SNSCreateTopicAPI defines the interface for the CreateTopic function.
// We use this interface to test the function using a mocked service.
type SNSTopicAPI interface {
	CreateTopic(context.Context, *sns.CreateTopicInput, ...func(*sns.Options)) (*sns.CreateTopicOutput, error)
	DeleteTopic(context.Context, *sns.DeleteTopicInput, ...func(*sns.Options)) (*sns.DeleteTopicOutput, error)
	ListTopics(context.Context, *sns.ListTopicsInput, ...func(*sns.Options)) (*sns.ListTopicsOutput, error)
}

func MakeTopic(c context.Context, api SNSTopicAPI, input *sns.CreateTopicInput) (*sns.CreateTopicOutput, error) {
	return api.CreateTopic(c, input)
}

func DestroyTopic(c context.Context, api SNSTopicAPI, input *sns.DeleteTopicInput) (*sns.DeleteTopicOutput, error) {
	return api.DeleteTopic(c, input)
}

func New() {
	ctx := context.Background()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	cfg, err := config.LoadDefaultConfig(ctx, func(o *config.LoadOptions) error {
		o.Region = "eu-west-2"
		o.HTTPClient = client
		return nil
	})
	if err != nil {
		panic(err)
	}

	snsClient := sns.NewFromConfig(cfg) // func(o *sns.Options) {
	// 	sns.WithEndpointResolver(sns.EndpointResolverFromURL("http://localstack.default.svc.cluster.local:4566"))
	// },

	out, err := MakeTopic(ctx, snsClient, &sns.CreateTopicInput{
		Name: aws.String("test-sns"),
	})
	if err != nil {
		panic(err)
	}

	log.Printf("%+v\n", *out.TopicArn)

	out2, err := snsClient.ListTopics(ctx, &sns.ListTopicsInput{})
	if err != nil {
		panic(err)
	}

	for _, t := range out2.Topics {
		fmt.Println(*t.TopicArn)
	}
}

func Delete() {
	ctx := context.Background()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	cfg, err := config.LoadDefaultConfig(ctx, func(o *config.LoadOptions) error {
		o.Region = "eu-west-2"
		o.HTTPClient = client
		return nil
	})
	if err != nil {
		panic(err)
	}

	snsClient := sns.NewFromConfig(cfg,
		func(o *sns.Options) {
			sns.WithEndpointResolver(sns.EndpointResolverFromURL("http://localstack.default.svc.cluster.local:4566"))
		},
	)

	out, err := DestroyTopic(ctx, snsClient, &sns.DeleteTopicInput{
		TopicArn: aws.String("arn:aws:sns:eu-west-2:735542962543:test-sns"),
	})
	if err != nil {
		panic(err)
	}

	log.Printf("%+v\n", out.ResultMetadata)
}
