package main

import (
	"context"
	"log"
	"os"

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

type NQObject struct {
	Spec struct {
		Name   string `json:"name,omitempty"`
		Region string `json:"region,omitempty"`
	} `json:"spec,omitempty"`
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

	notificationQueueInformer := cache.NewSharedIndexInformer(&cache.ListWatch{
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

	notificationQueueInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    AddResourceHandler,
		UpdateFunc: UpdateResourceHandler,
		DeleteFunc: DeleteResourceHandler,
	})

	stopCh := make(chan struct{})
	notificationQueueInformer.Run(stopCh)
	for range stopCh {
		os.Exit(0)
	}
}

func AddResourceHandler(obj any) {
	log.Printf("Hit add\n")
	nqo := unstructuredToNQ(obj)
	log.Printf("%+v\n", nqo.Spec.Name)
}
func UpdateResourceHandler(oldObj, newObj any) {
	log.Printf("Hit update\n")
}
func DeleteResourceHandler(obj any) {
	log.Printf("Hit delete\n")
	nqo := unstructuredToNQ(obj)
	log.Printf("%+v\n", nqo.Spec.Name)
}

func unstructuredToNQ(obj any) (nqo *NQObject) {
	nq, ok := obj.(*unstructured.Unstructured)
	if !ok {
		log.Printf("error in formatting of object\n")
	}

	err := runtime.DefaultUnstructuredConverter.FromUnstructured(nq.Object, nqo)
	if err != nil {
		log.Printf("error converting from unstructured to NQObject: %v\n", err)
	}
	return
}
