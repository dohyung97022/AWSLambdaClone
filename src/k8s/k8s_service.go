package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var Client *kubernetes.Clientset

func SetClient() error {
	var err error
	var config *rest.Config
	if config, err = rest.InClusterConfig(); err != nil {
		return err
	}
	if Client, err = kubernetes.NewForConfig(config); err != nil {
		return err
	}
	return nil
}
