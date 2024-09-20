package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/shresth72/kluster/controller"

	klient "github.com/shresth72/kluster/pkg/client/clientset/versioned"
	kinformerFactory "github.com/shresth72/kluster/pkg/client/informers/externalversions"
)

func main() {
	cluster := flag.String("cluster", "k3s", "Choose the cluster type: k3s or k8s")
	kubeconfig := flag.String("kubeconfig", "", "Location of the kubeconfig file")
	flag.Parse()

	if *kubeconfig == "" {
		if *cluster == "k8s" {
			*kubeconfig = "/home/shrestha/.kube/config"
		} else {
			*kubeconfig = "/etc/rancher/k3s/k3s.yaml"
		}
	}

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Printf("error building config from flags: %v\n", err)

		// If kubeconfig is not found, try to use in-cluster config
		config, err = rest.InClusterConfig()
		if err != nil {
			fmt.Printf("error fetching in-cluster config: %v\n", err)
			return
		}
	}

	klientset, err := klient.NewForConfig(config)
	if err != nil {
		log.Printf("getting klient set %s\n", err.Error())
	}

	stopCh := make(chan struct{})
	infoFactory := kinformerFactory.NewSharedInformerFactory(klientset, 20*time.Minute)

	c := controller.NewController(klientset, infoFactory.Shresth72().V1alpha1().Klusters())

	infoFactory.Start(stopCh)

	if err := c.Run(stopCh); err != nil {
		log.Printf("error running contoller: %v\n", err)
	}
}
