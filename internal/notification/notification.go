package notification

import (
	"context"
	"crypto/tls"
	"net/http"
	"strings"

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
func ListTopics(c context.Context, api SNSTopicAPI, input *sns.ListTopicsInput) (*sns.ListTopicsOutput, error) {
	return api.ListTopics(c, input)
}

func New(ctx context.Context, name, region string) *string {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	cfg, err := config.LoadDefaultConfig(ctx, func(o *config.LoadOptions) error {
		o.Region = region
		o.HTTPClient = client
		return nil
	})
	if err != nil {
		panic(err)
	}

	snsClient := sns.NewFromConfig(cfg)

	out, err := MakeTopic(ctx, snsClient, &sns.CreateTopicInput{
		Name: aws.String(name),
	})
	if err != nil {
		panic(err)
	}

	return out.TopicArn
}

func Delete(ctx context.Context, name, region string) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	cfg, err := config.LoadDefaultConfig(ctx, func(o *config.LoadOptions) error {
		o.Region = region
		o.HTTPClient = client
		return nil
	})
	if err != nil {
		panic(err)
	}
	snsClient := sns.NewFromConfig(cfg)

	out, err := ListTopics(ctx, snsClient, &sns.ListTopicsInput{})
	if err != nil {
		panic(err)
	}

	for _, topic := range out.Topics {
		arn := *topic.TopicArn

		splitTopic := strings.Split(arn, ":")
		topicName := splitTopic[len(splitTopic)-1]

		if topicName == name {
			_, err := DestroyTopic(ctx, snsClient, &sns.DeleteTopicInput{
				TopicArn: topic.TopicArn,
			})
			if err != nil {
				panic(err)
			}
			break
		}
	}
}
