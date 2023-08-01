package main

import (
	"context"
	"os"

	"github.com/izaakdale/simpleplane/internal/bucket"
	"github.com/kelseyhightower/envconfig"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type Specification struct {
	Group    string `envconfig:"GROUP"`
	Version  string `envconfig:"VERSION"`
	Resource string `envconfig:"RESOURCE"`
}

func main() {
	ctx := context.Background()

	var spec Specification
	envconfig.MustProcess("", &spec)

	gvr := schema.GroupVersionResource{
		Group:    spec.Group,
		Version:  spec.Version,
		Resource: spec.Resource,
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}
	dynCli := dynamic.NewForConfigOrDie(config)

	bucketInformer := cache.NewSharedIndexInformer(&cache.ListWatch{
		ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
			return dynCli.Resource(gvr).List(ctx, v1.ListOptions{})
		},
		WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
			return dynCli.Resource(gvr).Watch(ctx, v1.ListOptions{})
		},
	},
		&unstructured.Unstructured{},
		0,
		cache.Indexers{},
	)

	// hard coding since i am going to create roles etc for aws
	s3cli, err := bucket.NewClient(ctx, "eu-west-2")
	if err != nil {
		panic(err)
	}

	bucketInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    s3cli.AddResourceHandler,
		UpdateFunc: s3cli.UpdateResourceHandler,
		DeleteFunc: s3cli.DeleteResourceHandler,
	})

	stopCh := make(chan struct{})
	bucketInformer.Run(stopCh)
	for range stopCh {
		os.Exit(0)
	}
}
