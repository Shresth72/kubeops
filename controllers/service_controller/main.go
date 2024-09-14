package main

import (
	"flag"
	"fmt"
	"time"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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

	config.Timeout = 120 * time.Second

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("error building clientset: %v\n", err)
		return
	}

	stopCh := make(chan struct{})
	defer close(stopCh)

	informerFactory := informers.NewSharedInformerFactory(clientset, 10*time.Second)
	c := NewController(clientset, informerFactory.Apps().V1().Deployments())

	informerFactory.Start(stopCh)
	c.Run(stopCh)
}
