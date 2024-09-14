package main

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	appslisters "k8s.io/client-go/listers/apps/v1"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type Controller struct {
	clientSet     kubernetes.Interface
	depLister     appslisters.DeploymentLister
	depCacheSyncd cache.InformerSynced
	queue         workqueue.TypedRateLimitingInterface[any]
}

func NewController(
	clienset kubernetes.Interface,
	depInformer appsinformers.DeploymentInformer,
) *Controller {
	c := &Controller{
		clientSet:     clienset,
		depLister:     depInformer.Lister(),
		depCacheSyncd: depInformer.Informer().HasSynced,
		queue: workqueue.NewNamedRateLimitingQueue(
			workqueue.DefaultTypedControllerRateLimiter[any](),
			"ingress",
		),
	}

	depInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.handleAdd,
		DeleteFunc: c.handleDelete,
	})

	return c
}

func (c *Controller) Run(stopCh <-chan struct{}) {
	fmt.Println("Starting Controller...")

	if !cache.WaitForCacheSync(stopCh, c.depCacheSyncd) {
		fmt.Println("waiting for cache to be synced")
	}

	go wait.Until(c.worker, 1*time.Second, stopCh)

	<-stopCh
}

func (c *Controller) worker() {
	for c.processItems() {
	}
}

func (c *Controller) processItems() bool {
	item, shutdown := c.queue.Get()
	if shutdown {
		return false
	}
	defer c.queue.Forget(item)

	key, err := cache.MetaNamespaceKeyFunc(item)
	if err != nil {
		fmt.Printf("getting key from cache: %v\n", err)
		return false
	}

	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		fmt.Printf("splitting key into namespace and name: %v\n", err)
		return false
	}

	err = c.syncDeployment(ns, name)
	if err != nil {
		// retry
		c.queue.AddRateLimited(key)
		fmt.Printf("failed syncing deployment: %v; trying again\n", err)
		return true
	}

	return true
}

func (c *Controller) syncDeployment(ns, name string) error {
	ctx := context.Background()

	dep, err := c.depLister.Deployments(ns).Get(name)
	if err != nil {
		return fmt.Errorf("fetching deployment from informer: %v\n", err)
	}

	// create service
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dep.Name,
			Namespace: ns,
		},
		Spec: corev1.ServiceSpec{
			Selector: depLabels(*dep),
			Ports: []corev1.ServicePort{
				{
					Name: "http",
					Port: 80,
				},
			},
		},
	}

	_, err = c.clientSet.CoreV1().Services(ns).Create(ctx, &svc, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("creating service: %v\n", err)
	}

	// create ingress

	return nil
}

func (c *Controller) handleAdd(obj interface{}) {
	fmt.Println("\nDeployment Add called")
	c.queue.Add(obj)
}

func (c *Controller) handleDelete(obj interface{}) {
	fmt.Println("\nDeployment Delete called")
	c.queue.Add(obj)
}

func depLabels(dep appsv1.Deployment) map[string]string {
	return dep.Spec.Template.Labels
}
