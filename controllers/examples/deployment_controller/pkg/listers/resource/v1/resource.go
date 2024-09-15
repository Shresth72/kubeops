package v1

import (
	resourcev1 "github.com/Shresth72/deployment_controller/pkg/apis/resource/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/listers"
	"k8s.io/client-go/tools/cache"
)

// ResourceLister helps list Resources.
// All objects returned here must be treated as read-only.
type ResourceLister interface {
	// List lists all Resources in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*resourcev1.Resource, err error)
	// Resources returns an object that can list and get Resources.
	Resources(namespace string) ResourceNamespaceLister
	ResourceListerExpansion
}

// resourceLister implements the ResourceLister interface.
type resourceLister struct {
	listers.ResourceIndexer[*resourcev1.Resource]
}

// NewResourceLister returns a new ResourceLister.
func NewResourceLister(indexer cache.Indexer) ResourceLister {
	return &resourceLister{
		listers.New[*resourcev1.Resource](
			indexer,
			resourcev1.GetGroupResource("resource"),
		),
	}
}

// Resources returns an object that can list and get Resources.
func (s *resourceLister) Resources(namespace string) ResourceNamespaceLister {
	return resourceNamespaceLister{
		listers.NewNamespaced(s.ResourceIndexer, namespace),
	}
}

// ResourceNamespaceLister helps list and get Resources.
// All objects returned here must be treated as read-only.
type ResourceNamespaceLister interface {
	// List lists all Resources in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*resourcev1.Resource, err error)
	// Get retrieves the Resource from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*resourcev1.Resource, error)
	ResourceNamespaceListerExpansion
}

// resourceNamespaceLister implements the ResourceNamespaceLister
// interface.
type resourceNamespaceLister struct {
	listers.ResourceIndexer[*resourcev1.Resource]
}
