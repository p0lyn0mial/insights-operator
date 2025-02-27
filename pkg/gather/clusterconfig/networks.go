package clusterconfig

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	configv1client "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"

	"github.com/openshift/insights-operator/pkg/record"
)

// GatherClusterNetwork fetches the cluster Network - the Network with name cluster.
//
// The Kubernetes api https://github.com/openshift/client-go/blob/master/config/clientset/versioned/typed/config/v1/network.go#L50
// Response see https://docs.openshift.com/container-platform/4.3/rest_api/index.html#network-v1-config-openshift-io
//
// * Location in archive: config/network/
// * See: docs/insights-archive-sample/config/network
// * Id in config: networks
func GatherClusterNetwork(g *Gatherer, c chan<- gatherResult) {
	defer close(c)
	gatherConfigClient, err := configv1client.NewForConfig(g.gatherKubeConfig)
	if err != nil {
		c <- gatherResult{nil, []error{err}}
		return
	}
	records, errs := gatherClusterNetwork(g.ctx, gatherConfigClient)
	c <- gatherResult{records, errs}
}

func gatherClusterNetwork(ctx context.Context, configClient configv1client.ConfigV1Interface) ([]record.Record, []error) {
	config, err := configClient.Networks().Get(ctx, "cluster", metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, []error{err}
	}
	return []record.Record{{Name: "config/network", Item: record.JSONMarshaller{Object: config}}}, nil
}
