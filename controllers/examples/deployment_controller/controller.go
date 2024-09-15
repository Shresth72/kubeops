package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/time/rate"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/listers"
	appslisters "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"

	resourcev1 "github.com/Shresth72/deployment_controller/pkg/apis/resource/v1"
	// clientsetv1 "k8s.io/sample-controller/pkg/generated/clientset/versioned"
	// samplescheme "k8s.io/sample-controller/pkg/generated/clientset/versioned/scheme"
	// informers "k8s.io/sample-controller/pkg/generated/informers/externalversions/samplecontroller/v1alpha1"
	// listersv1 "k8s.io/sample-controller/pkg/generated/listers/samplecontroller/v1alpha1"
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
	sampleClientSet kubernetes.Interface

	deploymentsLister appslisters.DeploymentLister
	deploymentsSynced cache.InformerSynced
	resourceLister    listers.ResourceIndexer[*resourcev1.ResourceList]
	resourceSynced    cache.InformerSynced

	workqueue workqueue.TypedRateLimitingInterface[cache.ObjectName]
	recorder  record.EventRecorder
}

func NewController(
	ctx context.Context,
	kubeClientSet kubernetes.Interface,
	sampleClientSet kubernetes.Interface,
	deploymentInformer appsinformers.DeploymentInformer,
	resourceInformer informers.GenericInformer,
) *Controller {
	logger := klog.FromContext(ctx)
	utilruntime.Must(scheme.AddToScheme(scheme.Scheme))
	logger.V(4).Info("Creating event broadcaster")

	eventBroadcaster := record.NewBroadcaster(record.WithContext(ctx))
	eventBroadcaster.StartStructuredLogging(0)
	eventBroadcaster.StartRecordingToSink(
		&typedcorev1.EventSinkImpl{Interface: kubeClientSet.CoreV1().Events("")},
	)
	recorder := eventBroadcaster.NewRecorder(
		scheme.Scheme,
		corev1.EventSource{Component: controllerAgentName},
	)
	ratelimiter := workqueue.NewTypedMaxOfRateLimiter(
		workqueue.NewTypedItemExponentialFailureRateLimiter[cache.ObjectName](
			5*time.Millisecond,
			1000*time.Second,
		),
		&workqueue.TypedBucketRateLimiter[cache.ObjectName]{
			Limiter: rate.NewLimiter(rate.Limit(50), 300),
		},
	)

	controller := &Controller{
		kubeClientSet:     kubeClientSet,
		sampleClientSet:   sampleClientSet,
		deploymentsLister: deploymentInformer.Lister(),
		deploymentsSynced: deploymentInformer.Informer().HasSynced,
		resourceLister:    resourceInformer.Lister(),
		resourceSynced:    resourceInformer.Informer().HasSynced,
		workqueue:         workqueue.NewTypedRateLimitingQueue(ratelimiter),
		recorder:          recorder,
	}
}
