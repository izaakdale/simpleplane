package bucket

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	api S3API
}

func NewClient(ctx context.Context, region string) (*Client, error) {
	var cli Client

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
		return nil, err
	}

	cli.api = s3.NewFromConfig(cfg)

	return &cli, nil
}
