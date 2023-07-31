package main

import (
	"context"
	"log"
	"os"

	"github.com/izaakdale/simpleplane/internal/notification"
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
	Group     string   `envconfig:"GROUP"`
	Version   string   `envconfig:"VERSION"`
	Resources []string `envconfig:"RESOURCE"`
}

func main() {
	ctx := context.Background()

	var spec Specification
	envconfig.MustProcess("", &spec)

	log.Printf("%+v\n", spec.Resources)

	gvr := schema.GroupVersionResource{
		Group:    spec.Group,
		Version:  spec.Version,
		Resource: spec.Resources[0],
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}
	dynCli := dynamic.NewForConfigOrDie(config)

	notificationInformer := cache.NewSharedIndexInformer(&cache.ListWatch{
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

	notificationInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    notification.AddResourceHandler,
		UpdateFunc: notification.UpdateResourceHandler,
		DeleteFunc: notification.DeleteResourceHandler,
	})

	stopCh := make(chan struct{})
	notificationInformer.Run(stopCh)
	for range stopCh {
		os.Exit(0)
	}
}
