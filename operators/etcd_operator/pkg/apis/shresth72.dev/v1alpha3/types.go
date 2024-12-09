package v1alpha3

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

const ()

type EtcdadmClusterSpec struct {
	Replicas               *int32                 `json:"replicas,omitempty"`
	InfrastructureTemplate corev1.ObjectReference `json:"infrastructureTemplate"`
}

type EtcdadmClusterStatus struct {
	ReadyReplicas      int32                `json:"replicas,omitempty"`
	InitMachineAddress string               `json:"initMachineAddress"`
	Initialized        bool                 `json:"initialized"`
	Ready              bool                 `json:"reader"`
	CreationComplete   bool                 `json:"creationcomplete"`
	Endpoints          string               `json:"endpoints"`
	Selector           string               `json:"selector,omitempty"`
	ObservedGeneration int64                `json:"observedGeneration,omitempty"`
	Conditions         clusterv1.Conditions `json:"conditions,omitempty"`
}

type EtcdadmCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EtcdadmClusterSpec   `json:"spec,omitempty"`
	Status EtcdadmClusterStatus `json:"status,omitempty"`
}

func (in *EtcdadmCluster) GetConditions() clusterv1.Conditions {
	return in.Status.Conditions
}

func (in *EtcdadmCluster) SetConditions(conditions clusterv1.Conditions) {
	in.Status.Conditions = conditions
}

type EtcdadmClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []EtcdadmCluster `json:"items"`
}
