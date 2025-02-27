package clusterconfig

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	configv1client "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"

	"github.com/openshift/insights-operator/pkg/record"
)

// GatherClusterOAuth fetches the cluster OAuth - the OAuth with name cluster.
//
// The Kubernetes api https://github.com/openshift/client-go/blob/master/config/clientset/versioned/typed/config/v1/oauth.go#L50
// Response see https://docs.openshift.com/container-platform/4.3/rest_api/index.html#oauth-v1-config-openshift-io
//
// * Location in archive: config/oauth/
// * See: docs/insights-archive-sample/config/oauth
// * Id in config: oauths
func GatherClusterOAuth(g *Gatherer, c chan<- gatherResult) {
	defer close(c)
	gatherConfigClient, err := configv1client.NewForConfig(g.gatherKubeConfig)
	if err != nil {
		c <- gatherResult{nil, []error{err}}
		return
	}
	records, errs := gatherClusterOAuth(g.ctx, gatherConfigClient)
	c <- gatherResult{records, errs}
}

func gatherClusterOAuth(ctx context.Context, configClient configv1client.ConfigV1Interface) ([]record.Record, []error) {
	config, err := configClient.OAuths().Get(ctx, "cluster", metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, []error{err}
	}
	return []record.Record{{Name: "config/oauth", Item: record.JSONMarshaller{Object: config}}}, nil
}
