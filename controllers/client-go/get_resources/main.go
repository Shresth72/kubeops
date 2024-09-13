package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Get all pods. deployments from a kubernetes cluster
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

	ctx := context.Background()
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

	// Set custom timeout for api server to respond
	config.Timeout = 120 * time.Second

	// Pods are in CoreV1 api
	// Pod implements runtime.Object (Kubernetes object)
	// 1. TypeMeta: impl SetGroupVersionKind, GroupVersionKind
	//    - apiVersion: apps/v1
	//    - kind: Deployment
	//  - impl DeepCopyObject
	//
	// 2. ObjectMeta
	// 3. Spec
	// 4. Status

	// ClientSet contains the clients for groups. Each group has exactly one version included in a ClientSet

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("error building clientset: %v\n", err)
		return
	}

	pods, err := clientset.CoreV1().Pods("default").List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Printf("error fetching pods: %v\n", err)
		return
	}

	fmt.Println("Pods from default namespace")
	for _, pod := range pods.Items {
		fmt.Printf("%s\n", pod.Name)
	}

	// Deployments are in appv1 api
	deployments, err := clientset.AppsV1().Deployments("default").List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Printf("error fetching deployments: %v\n", err)
		return
	}

	fmt.Println("\nDeployments from default namespace")
	for _, d := range deployments.Items {
		fmt.Printf("%s\n", d.Name)
	}
}
