package main

import (
	"context"
	"fmt"
	"time"

	resourcev1 "github.com/Shresth72/deployment_controller/pkg/apis/resource/v1"

	"golang.org/x/time/rate"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	appslisters "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"

	clientset "k8s.io/sample-controller/pkg/generated/clientset/versioned"
	samplescheme "k8s.io/sample-controller/pkg/generated/clientset/versioned/scheme"
	informers "k8s.io/sample-controller/pkg/generated/informers/externalversions/samplecontroller/v1alpha1"
	listers "k8s.io/sample-controller/pkg/generated/listers/samplecontroller/v1alpha1"
)

const controllerAgentName = "deployment-controller"

var resourceAgentName = "foo"

const (
	SuccessSynced     = "Synced"
	ErrResourceExists = "ErrResourceExists"
	FieldManager      = controllerAgentName
)

var (
	MessageResourceExists = fmt.Sprintf(
		"Resource %q already exists and is not managed by %s",
		resourceAgentName,
	)
	MessageResourceSynced = fmt.Sprintf("%s synced successfully", resourceAgentName)
)

type Controller struct {
	kubeClientSet kubernetes.Interface
	// For our own API group
	sampleClientSet clientset.Interface

	deploymentsLister appslisters.DeploymentLister
	deploymentsSynced cache.InformerSynced
	foosLister        listers.FooLister
	foosSynced        cache.InformerSynced

	workqueue workqueue.TypedRateLimitingInterface[cache.ObjectName]
	recorder  record.EventRecorder
}

func NewController(
	ctx context.Context,
	kubeClientSet kubernetes.Interface,
	sampleClientSet clientset.Interface,
	deploymentInformer appsinformers.DeploymentInformer,
	fooInformer informers.FooInformer,
)
