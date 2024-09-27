package controller

import (
	"context"
	"log"
	"time"

	"github.com/kanisterio/kanister/pkg/poll"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"

	"github.com/shresth72/kluster/pkg/do"

	v1alpha1 "github.com/shresth72/kluster/pkg/apis/shresth72.dev/v1alpha1"
	klientset "github.com/shresth72/kluster/pkg/client/clientset/versioned"
	customscheme "github.com/shresth72/kluster/pkg/client/clientset/versioned/scheme"
	kinformer "github.com/shresth72/kluster/pkg/client/informers/externalversions/shresth72.dev/v1alpha1"
	klister "github.com/shresth72/kluster/pkg/client/listers/shresth72.dev/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type Controller struct {
	client kubernetes.Interface

	// clientset for custom resource kluster
	klient klientset.Interface
	// kluster has synced
	klusterSynced cache.InformerSynced
	// lister
	kLister klister.KlusterLister
	// queue
	queue workqueue.TypedRateLimitingInterface[any]
	// event recorder
	recorder record.EventRecorder
}

func NewController(
	client kubernetes.Interface,
	klient klientset.Interface,
	klusterInformer kinformer.KlusterInformer,
) *Controller {
	runtime.Must(customscheme.AddToScheme(scheme.Scheme))

	eventBroadCaster := record.NewBroadcaster()
	eventBroadCaster.StartStructuredLogging(0)

	eventBroadCaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{
		Interface: client.CoreV1().Events(""),
	})
	recorder := eventBroadCaster.NewRecorder(
		scheme.Scheme,
		corev1.EventSource{Component: "Kluster"},
	)

	c := &Controller{
		client:        client,
		klient:        klient,
		klusterSynced: klusterInformer.Informer().HasSynced,
		kLister:       klusterInformer.Lister(),
		queue: workqueue.NewNamedRateLimitingQueue(
			workqueue.DefaultTypedControllerRateLimiter[any](),
			"kluster",
		),
		recorder: recorder,
	}

	klusterInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.handleAdd,
		DeleteFunc: c.handleDelete,
	})

	return c
}

func (c *Controller) Run(stopCh chan struct{}) error {
	if ok := cache.WaitForCacheSync(stopCh, c.klusterSynced); !ok {
		log.Println("cache was not synced")
	}

	go wait.Until(c.worker, time.Second, stopCh)

	// Block until it recieves a signal
	<-stopCh
	return nil
}

func (c *Controller) worker() {
	for c.processNextItem() {
	}
}

func (c *Controller) processNextItem() bool {
	item, shutDown := c.queue.Get()
	if shutDown {
		log.Println("process shutdown")
		return false
	}
	defer c.queue.Forget(item)

	key, err := cache.MetaNamespaceKeyFunc(item)
	if err != nil {
		log.Printf("error calling Namespace key func: %v on cache for item", err)
		return false
	}

	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		log.Printf("error splitting meta namespace: %v", err)
		return false
	}

	kluster, err := c.kLister.Klusters(ns).Get(name)
	if err != nil {
		log.Printf("error getting klusters resource from lister in namespace '%s':  %v", ns, err)
		return false
	}

	log.Printf("current kluster spec: %+v\n", kluster.Spec)

	// Delete the cluster from DO if also deleted from local cluster
	_, err = c.klient.Shresth72V1alpha1().
		Klusters(ns).
		Get(context.Background(), name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		log.Printf("Deployment '%s' was deleted\n", name)

		err := do.Delete(c.client, kluster.Spec)
		if err != nil {
			log.Printf("failed to delete '%s' kluster\n", name)
			return false
		}

		log.Printf("Kluster '%s' on Digital Ocean successfully deleted", name)
		return true
	}

	// Create and manage the kluster of Digital Ocean
	// Persisting ClusterID for easy deletion
	clusterId, err := do.Create(c.client, kluster.Spec)
	if err != nil {
		log.Printf("error creating cluster on digital ocean: %v", err)
		return false
	}

	c.recorder.Event(
		kluster,
		corev1.EventTypeNormal,
		"ClusterCreation",
		"Digital Ocean API was called to create the cluster",
	)

	log.Printf("The Digital Ocean ClusterId that we have is: %s\n", clusterId)

	// TODO: check if returning false is right here

	if err = c.updateStatus(clusterId, "creating", kluster); err != nil {
		log.Printf("error updating '%s' cluster status: %v", kluster.Name, err)
		return false
	}

	// query Digital Ocean API to make sure cluster state is running
	if err = c.waitForCluster(kluster.Spec, clusterId); err != nil {
		log.Printf("error waiting for cluster '%s' to be running: %v", kluster.Name, err)
		return false
	}

	if err = c.updateStatus(clusterId, "running", kluster); err != nil {
		log.Printf(
			"error updating '%s' cluster status after waiting for cluster: %v",
			kluster.Name,
			err,
		)
		return false
	}

	c.recorder.Event(
		kluster,
		corev1.EventTypeNormal,
		"ClusterCreationCompleted",
		"Digital Ocean Cluster creation was completed",
	)

	return true
}

func (c *Controller) waitForCluster(spec v1alpha1.KlusterSpec, clusterId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	return poll.Wait(ctx, func(ctx context.Context) (bool, error) {
		state, err := do.ClusterState(c.client, spec, clusterId)
		if err != nil {
			return false, err
		}

		if state == "running" {
			return true, nil
		}

		return false, nil
	})
}

func (c *Controller) updateStatus(id, progress string, kluster *v1alpha1.Kluster) error {
	// Fetch the latest changes
	k, err := c.klient.Shresth72V1alpha1().
		Klusters(kluster.Namespace).
		Get(context.Background(), kluster.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	k.Status.KlusterId = id
	k.Status.Progress = progress

	_, err = c.klient.Shresth72V1alpha1().
		Klusters(kluster.Namespace).
		UpdateStatus(context.Background(), k, metav1.UpdateOptions{})
	return err
}

func (c *Controller) handleAdd(obj interface{}) {
	log.Println("handleAdd was called")
	c.queue.Add(obj)
}

func (c *Controller) handleDelete(obj interface{}) {
	log.Println("handleDelete was called")
	c.queue.Add(obj)
}
