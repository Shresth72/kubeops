// Get all pods, deployments from a kubernetes cluster

package main

import (
	"flag"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

/*
 Pods are in CoreV1 api, Deployments are in AppsV1 api
 Pods, Deployments, etc implement runtime.Object (Kubernetes object)

 1. TypeMeta: impl SetGroupVersionKind, GroupVersionKind
    - apiVersion: apps/v1 (resource)
    - kind: Deployment | Pod | Service
        - RestMapper: Pass GroupVersionResource to get the GroupVersionKind
        - ObjectKinds: Pass k8s object to get the GroupVersionKind (if the object is already registered through AddKnownTypes)

 2. ObjectMeta
 3. Spec
 4. Status

 - impl DeepCopyObject
*/

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

	// Set custom timeout for api server to respond
	config.Timeout = 120 * time.Second

	// ClientSet contains the clients for groups. Each group has exactly one version included in a ClientSet
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("error building clientset: %v\n", err)
		return
	}

	// Using Shared InformerFactory for caching resource information from the API server. Hence, avoiding repetitive load on the api server
	// Instead of calling List()
	// pods, err := clientset.CoreV1().Pods("default").List(ctx, metav1.ListOptions{})
	// deployments, err := clientset.AppsV1().Deployments("default").List(ctx, metav1.ListOptions{})

	informerFactor := informers.NewSharedInformerFactory(clientset, 10*time.Minute)

	// For custom informers from chosen namespace
	_ = informers.NewFilteredSharedInformerFactory(
		clientset,
		10*time.Minute,
		"test",
		func(lo *metav1.ListOptions) {
			lo.LabelSelector = ""
			lo.APIVersion = ""
		},
	)

	podInformer := informerFactor.Core().V1().Pods()

	// For handling events, queues should be used to handle each event in a seperate go routine
	// Allowing for rollbacks and better workflow
	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			fmt.Println("add was called")
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			// We can compare newObj resourceVersion to check if it was updated in case of ressync
			fmt.Println("updated was called")
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println("delete was called")
		},
	})

	informerFactor.Start(wait.NeverStop)
	informerFactor.WaitForCacheSync(wait.NeverStop)

	// Also, for changing resource configuration,
	// DeepCopyObject should be used to prevent cache consistency issues
	pod, err := podInformer.Lister().Pods("default").Get("default")
	fmt.Println(pod.Name)
}
