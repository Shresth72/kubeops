package do

import (
	"context"
	"fmt"
	"strings"

	"k8s.io/client-go/kubernetes"

	"github.com/digitalocean/godo"

	v1alpha1 "github.com/shresth72/kluster/pkg/apis/shresth72.dev/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Create(client kubernetes.Interface, spec v1alpha1.KlusterSpec) (string, error) {
	token, err := getToken(client, spec.TokenSecret)
	if err != nil {
		return "", err
	}

	tokenClient := godo.NewFromToken(token)
	fmt.Println(tokenClient)

	// TODO: Validate if the spec had all the required fileds
	request := &godo.KubernetesClusterCreateRequest{
		Name:        spec.Name,
		RegionSlug:  spec.Region,
		VersionSlug: spec.Version,
		NodePools: []*godo.KubernetesNodePoolCreateRequest{
			{
				Size:  spec.NodePools[0].Size,
				Name:  spec.NodePools[0].Name,
				Count: spec.NodePools[0].Count,
			},
		},
	}

	cluster, _, err := tokenClient.Kubernetes.Create(context.Background(), request)
	if err != nil {
		return "", err
	}

	return cluster.ID, nil
}

func ClusterState(c kubernetes.Interface, spec v1alpha1.KlusterSpec, id string) (string, error) {
	token, err := getToken(c, spec.TokenSecret)
	if err != nil {
		return "", err
	}

	client := godo.NewFromToken(token)
	cluster, _, err := client.Kubernetes.Get(context.Background(), id)

	return string(cluster.Status.State), err
}

func getToken(client kubernetes.Interface, secret string) (string, error) {
	namespace := strings.Split(secret, "/")[0]
	name := strings.Split(secret, "/")[1]

	s, err := client.CoreV1().
		Secrets(namespace).
		Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	return string(s.Data["token"]), nil
}
