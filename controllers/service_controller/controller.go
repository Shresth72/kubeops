package main

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
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

	// check if the object has been deleted from k8s cluster
	ctx := context.Background()
	_, err = c.clientSet.AppsV1().Deployments(ns).Get(ctx, name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		fmt.Printf("Deployment '%s' was deleted\n", name)

		// delete service
		err := c.clientSet.CoreV1().Secrets(ns).Delete(ctx, name, metav1.DeleteOptions{})
		if err != nil {
			fmt.Printf("failed to delete '%s' service\n", name)
			return false
		}

		// delete ingress
		err = c.clientSet.NetworkingV1().Ingresses(ns).Delete(ctx, name, metav1.DeleteOptions{})
		if err != nil {
			fmt.Printf("failed to delete '%s' ingress resource\n", name)
			return false
		}

		return true
	}

	err = c.syncDeployment(ns, name)
	if err != nil {
		// retry
		c.queue.AddRateLimited(key)
		fmt.Printf("error syncing deployment: %v; trying again\n", err)
		return false
	}

	return true
}

func (c *Controller) syncDeployment(ns, name string) error {
	ctx := context.Background()

	dep, err := c.depLister.Deployments(ns).Get(name)
	if err != nil {
		return fmt.Errorf("error fetching deployment from informer: %v\n", err)
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

	service, err := c.clientSet.CoreV1().Services(ns).Create(ctx, &svc, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("error creating service: %v\n", err)
	}

	// create ingress resource (/path)
	return createIngress(ctx, c.clientSet, service)
}

func createIngress(ctx context.Context, client kubernetes.Interface, svc *corev1.Service) error {
	pathType := "Prefix"

	ingress := netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      svc.Name,
			Namespace: svc.Namespace,
			Annotations: map[string]string{
				"nginx.ingress.kubernetes.io/rewrite-target": "/",
			},
		},
		Spec: netv1.IngressSpec{
			Rules: []netv1.IngressRule{
				{
					IngressRuleValue: netv1.IngressRuleValue{
						HTTP: &netv1.HTTPIngressRuleValue{
							Paths: []netv1.HTTPIngressPath{
								{
									Path:     fmt.Sprintf("%s", svc.Name),
									PathType: (*netv1.PathType)(&pathType),
									Backend: netv1.IngressBackend{
										Service: &netv1.IngressServiceBackend{
											Name: svc.Name,
											Port: netv1.ServiceBackendPort{
												Number: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := client.NetworkingV1().
		Ingresses(svc.Namespace).
		Create(ctx, &ingress, metav1.CreateOptions{})

	return err
}

func (c *Controller) handleAdd(obj interface{}) {
	dep, ok := obj.(*appsv1.Deployment)
	if !ok {
		fmt.Println("unexpected object type")
	}

	ctx := context.Background()
	_, err := c.clientSet.CoreV1().Services(dep.Namespace).Get(ctx, dep.Name, metav1.GetOptions{})
	if err == nil {
		fmt.Printf("Deployment '%s' already exists\n", dep.Name)
		return
	}

	fmt.Printf("\nDeployment Add called for '%s'", dep.Name)
	c.queue.Add(obj)
}

func (c *Controller) handleDelete(obj interface{}) {
	fmt.Println("\nDeployment Delete called")
	c.queue.Add(obj)
}

func depLabels(dep appsv1.Deployment) map[string]string {
	return dep.Spec.Template.Labels
}
