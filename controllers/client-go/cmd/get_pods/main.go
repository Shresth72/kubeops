package main

import (
	"context"
	"flag"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Get all pods from a kubernetes cluster
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
		fmt.Printf("error building config: %v\n", err)
		return
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("error building clientset: %v\n", err)
		return
	}

	ctx := context.Background()

	// Pods are in corev1 api
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

	fmt.Println("\nDeployments from default namespace")
	for _, d := range deployments.Items {
		fmt.Printf("%s\n", d.Name)
	}
}
