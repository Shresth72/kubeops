package controller

import (
	"log"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	"github.com/shresth72/kluster/pkg/do"

	klientset "github.com/shresth72/kluster/pkg/client/clientset/versioned"
	kinformer "github.com/shresth72/kluster/pkg/client/informers/externalversions/shresth72.dev/v1alpha1"
	klister "github.com/shresth72/kluster/pkg/client/listers/shresth72.dev/v1alpha1"
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
}

func NewController(
	client kubernetes.Interface,
	klient klientset.Interface,
	klusterInformer kinformer.KlusterInformer,
) *Controller {
	c := &Controller{
		client:        client,
		klient:        klient,
		klusterSynced: klusterInformer.Informer().HasSynced,
		kLister:       klusterInformer.Lister(),
		queue: workqueue.NewNamedRateLimitingQueue(
			workqueue.DefaultTypedControllerRateLimiter[any](),
			"kluster",
		),
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

	// Create and manage the kluster of Digital Ocean
	clusterId, err := do.Create(c.client, kluster.Spec)
	if err != nil {
		log.Printf("error creating cluster on digital ocean: %v", err)
		return false
	}

	log.Printf("The Digital Ocean ClusterId that we have is: %s\n", clusterId)

	return true
}

func (c *Controller) handleAdd(obj interface{}) {
	log.Println("handleAdd was called")
	c.queue.Add(obj)
}

func (c *Controller) handleDelete(obj interface{}) {
	log.Println("handleDelete was called")
	c.queue.Add(obj)
}
