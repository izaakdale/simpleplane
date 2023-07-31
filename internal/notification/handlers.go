package notification

import (
	"context"
	"log"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

type NotificationObject struct {
	Spec struct {
		Name   string `json:"name,omitempty"`
		Region string `json:"region,omitempty"`
	} `json:"spec,omitempty"`
}

func AddResourceHandler(obj any) {
	log.Printf("Hit add\n")

	nq, ok := obj.(*unstructured.Unstructured)
	if !ok {
		log.Printf("error in formatting of object\n")
	}

	var no NotificationObject
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(nq.Object, &no)
	if err != nil {
		log.Printf("error converting from unstructured to NQObject: %v\n", err)
	}

	log.Printf("%s : %s\n", no.Spec.Name, no.Spec.Region)

	_ = New(context.TODO(), no.Spec.Name, no.Spec.Region)
}

func UpdateResourceHandler(oldObj, newObj any) {
	log.Printf("Hit update\n")
}

func DeleteResourceHandler(obj any) {
	log.Printf("Hit delete\n")

	nq, ok := obj.(*unstructured.Unstructured)
	if !ok {
		log.Printf("error in formatting of object\n")
	}

	var no NotificationObject
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(nq.Object, &no)
	if err != nil {
		log.Printf("error converting from unstructured to NQObject: %v\n", err)
	}
	Delete(context.TODO(), no.Spec.Name, no.Spec.Region)
}
