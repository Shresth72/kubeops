package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	klient "github.com/shresth72/kluster/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	klusters, err := klientset.Shresth72V1alpha1().
		Klusters("").
		List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Printf("listing klusters %s\n", err.Error())
	}

	fmt.Printf("length of klusters is: %d\n", len(klusters.Items))
}
