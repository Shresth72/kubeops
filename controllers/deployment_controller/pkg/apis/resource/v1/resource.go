package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// impl TypeMeta, ObjectMeta, Spec, Status for Kubernetes Object

type Resource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResourceSpec   `json:"spec"`
	Status ResourceStatus `json:"status"`
}

type ResourceSpec struct {
	DeploymentName string `json:"deploymentName"`
	Replicas       *int32 `json:"replicas"`
}

type ResourceStatus struct {
	AvailableReplicas int32 `json:"availableReplicas"`
}

type ResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Resource `json:"items"`
}
