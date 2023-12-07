package bucket

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/kelseyhightower/envconfig"
)

type Client struct {
	api S3API
}

type provider struct {
	AccessKeyID     string `envconfig:"AWS_ACCESS_KEY_ID"`
	SecretAccessKey string `envconfig:"AWS_SECRET_ACCESS_KEY"`
}

func (p *provider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	err := envconfig.Process("", p)
	if err != nil {
		return aws.Credentials{}, err
	}

	return aws.Credentials{
		AccessKeyID:     p.AccessKeyID,
		SecretAccessKey: p.SecretAccessKey,
	}, nil
}

func NewClient(ctx context.Context, region string) (*Client, error) {
	var cli Client

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(&provider{}), func(o *config.LoadOptions) error {
		o.Region = region
		o.HTTPClient = client
		return nil
	})
	if err != nil {
		return nil, err
	}

	cli.api = s3.NewFromConfig(cfg)

	return &cli, nil
}
