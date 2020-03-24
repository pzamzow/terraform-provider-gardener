package shoot

import (
	"encoding/json"
	"testing"

	gcpAlpha1 "github.com/gardener/gardener-extension-provider-gcp/pkg/apis/gcp/v1alpha1"
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	cmp "github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/schema"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func getProviderGCPControlPlaneConfigTestData() (map[string]interface{}, *corev1beta1.ProviderConfig) {
	gcpConfig, _ := json.Marshal(gcpAlpha1.ControlPlaneConfig{
		TypeMeta: v1.TypeMeta{
			APIVersion: "gcp.provider.extensions.gardener.cloud/v1alpha1",
			Kind:       "ControlPlaneConfig",
		},
		Zone: "zone1",
	})
	gcpRaw := map[string]interface{}{
		"zone": "zone1",
	}
	gcp := &corev1beta1.ProviderConfig{
		RawExtension: runtime.RawExtension{
			Raw: gcpConfig,
		},
	}

	return gcpRaw, gcp
}

func getProviderGCPInfrastructureConfigTestData() (map[string]interface{}, *corev1beta1.ProviderConfig) {
	internal := "test"
	fooCloudNat := int32(2)
	aggregationInterval := "foo"
	flowSampling := float32(0.5)
	metadata := "bar"
	gcpConfig, _ := json.Marshal(gcpAlpha1.InfrastructureConfig{
		TypeMeta: v1.TypeMeta{
			APIVersion: "gcp.provider.extensions.gardener.cloud/v1alpha1",
			Kind:       "InfrastructureConfig",
		},
		Networks: gcpAlpha1.NetworkConfig{
			VPC: &gcpAlpha1.VPC{
				Name: "foo", CloudRouter: &gcpAlpha1.CloudRouter{
					Name: "bar",
				},
			},
			CloudNAT: &gcpAlpha1.CloudNAT{
				MinPortsPerVM: &fooCloudNat,
			},
			Internal: &internal,
			Worker:   "",
			Workers:  "10.250.0.0/19",
			FlowLogs: &gcpAlpha1.FlowLogs{
				AggregationInterval: &aggregationInterval,
				FlowSampling:        &flowSampling,
				Metadata:            &metadata,
			},
		},
	})
	gcpRaw := map[string]interface{}{
		"networks": []interface{}{
			map[string]interface{}{
				"vpc": []interface{}{
					map[string]interface{}{
						"name": "foo",
						"cloud_router": []interface{}{
							map[string]interface{}{
								"name": "bar",
							},
						},
					},
				},
				"cloud_nat": []interface{}{
					map[string]interface{}{
						"min_ports_per_vm": 2,
					},
				},
				"internal": "test",
				"workers":  "10.250.0.0/19",
				"flow_logs": []interface{}{
					map[string]interface{}{
						"aggregation_interval": aggregationInterval,
						"flow_sampling":        0.5,
						"metadata":             metadata,
					},
				},
			},
		},
	}
	gcp := &corev1beta1.ProviderConfig{
		RawExtension: runtime.RawExtension{
			Raw: gcpConfig,
		},
	}

	return gcpRaw, gcp
}

func TestExpandProviderGCPInfrastructureConfig(t *testing.T) {
	gcpRaw, gcp := getProviderGCPInfrastructureConfigTestData()
	data := schema.TestResourceDataRaw(t, GCPInfrastructureConfigResource().Schema, gcpRaw)
	out := ExpandGCPInfrastructureConfig([]interface{}{data.Get("")})
	if diff := cmp.Diff(gcp, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}

func TestFlattenProviderGCPInfrastructureConfig(t *testing.T) {
	gcpRaw, gcp := getProviderGCPInfrastructureConfigTestData()
	out := FlattenGCPInfrastructureConfig(gcp)[0].(map[string]interface{})
	data := schema.TestResourceDataRaw(t, GCPInfrastructureConfigResource().Schema, gcpRaw).Get("")
	if diff := cmp.Diff(data, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}

func TestExpandProviderGCPControlPlaneConfig(t *testing.T) {
	gcpRaw, gcp := getProviderGCPControlPlaneConfigTestData()
	data := schema.TestResourceDataRaw(t, GCPControlPlaneConfigResource().Schema, gcpRaw)
	out := ExpandGCPControlPlaneConfig([]interface{}{data.Get("")})
	if diff := cmp.Diff(gcp, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}

func TestFlattenProviderGCPControlPlaneConfig(t *testing.T) {
	gcpRaw, gcp := getProviderGCPControlPlaneConfigTestData()
	out := FlattenGCPControlPlaneConfig(gcp)[0].(map[string]interface{})
	data := schema.TestResourceDataRaw(t, GCPControlPlaneConfigResource().Schema, gcpRaw).Get("")
	if diff := cmp.Diff(data, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}
