package main

import (
	"flag"
	"time"

	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"k8s.io/sample-controller/pkg/signals"
	// gcp "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	clientset "k8s.io/sample-controller/pkg/generated/clientset/versioned"
	informers "k8s.io/sample-controller/pkg/generated/informers/externalversions"
)
