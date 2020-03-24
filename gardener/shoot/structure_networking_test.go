package shoot

import (
	"testing"

	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	cmp "github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/schema"
)

func getNetworkingTestData() (map[string]interface{}, *corev1beta1.Networking) {
	nodes := "10.250.0.0/19"
	pods := "100.96.0.0/11"
	services := "100.64.0.0/13"
	networkingRaw := map[string]interface{}{
		"nodes":    nodes,
		"pods":     pods,
		"services": services,
		"type":     "calico",
	}
	networking := &corev1beta1.Networking{
		Nodes:    &nodes,
		Pods:     &pods,
		Services: &services,
		Type:     "calico",
	}

	return networkingRaw, networking
}

func TestExpandNetworking(t *testing.T) {
	networkingRaw, networking := getNetworkingTestData()
	data := schema.TestResourceDataRaw(t, NetworkingResource().Schema, networkingRaw)
	out := ExpandNetworking([]interface{}{data.Get("")})
	if diff := cmp.Diff(networking, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}

func TestFlattenNetworking(t *testing.T) {
	networkingRaw, networking := getNetworkingTestData()
	out := FlattenNetworking(networking)[0].(map[string]interface{})
	data := schema.TestResourceDataRaw(t, NetworkingResource().Schema, networkingRaw).Get("")
	if diff := cmp.Diff(data, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}
