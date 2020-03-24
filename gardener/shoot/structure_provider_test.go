package shoot

import (
	"testing"

	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	cmp "github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/schema"
)

func getProviderTestData(name string) (map[string]interface{}, *corev1beta1.Provider) {
	workerRaw, worker := getWorkerTestData()

	var providerInfrastructureConfig *corev1beta1.ProviderConfig
	var providerInfrastructureConfigRaw map[string]interface{}
	var providerControlPlaneConfig *corev1beta1.ProviderConfig
	var providerControlPlaneConfigRaw []interface{}
	switch name {
	case "azure":
		azureRaw, azure := getProviderAzureInfrastructureConfigTestData()
		providerInfrastructureConfigRaw = azureRaw
		providerInfrastructureConfig = azure
	case "aws":
		awsRaw, aws := getProviderAWSInfrastructureConfigTestData()
		providerInfrastructureConfigRaw = awsRaw
		providerInfrastructureConfig = aws
	case "gcp":
		gcpInfRaw, gcpInf := getProviderGCPInfrastructureConfigTestData()
		gcpCtlRaw, gcpCtl := getProviderGCPControlPlaneConfigTestData()
		providerInfrastructureConfigRaw = gcpInfRaw
		providerInfrastructureConfig = gcpInf
		providerControlPlaneConfigRaw = []interface{}{
			map[string]interface{}{
				name: []interface{}{gcpCtlRaw},
			},
		}
		providerControlPlaneConfig = gcpCtl
	}

	providerRaw := map[string]interface{}{
		"type": name,
		"infrastructure_config": []interface{}{
			map[string]interface{}{
				name: []interface{}{providerInfrastructureConfigRaw},
			},
		},
		"control_plane_config": providerControlPlaneConfigRaw,
		"worker":               []interface{}{workerRaw},
	}
	provider := &corev1beta1.Provider{
		Type:                 name,
		InfrastructureConfig: providerInfrastructureConfig,
		ControlPlaneConfig:   providerControlPlaneConfig,
		Workers:              []corev1beta1.Worker{*worker},
	}

	return providerRaw, provider
}

func TestExpandProvider(t *testing.T) {
	testCases := []string{"aws", "azure", "gcp"}

	for _, tc := range testCases {
		providerRaw, provider := getProviderTestData(tc)
		data := schema.TestResourceDataRaw(t, ProviderResource().Schema, providerRaw)
		out := ExpandProvider([]interface{}{data.Get("")})
		if diff := cmp.Diff(provider, out); diff != "" {
			t.Fatalf("Error matching output and expected: \n%s", diff)
		}
	}
}

func TestFlattenProvider(t *testing.T) {
	testCases := []string{"aws", "azure", "gcp"}

	for _, tc := range testCases {
		providerRaw, provider := getProviderTestData(tc)
		out := FlattenProvider(provider)[0].(map[string]interface{})
		data := schema.TestResourceDataRaw(t, ProviderResource().Schema, providerRaw).Get("")
		if diff := cmp.Diff(data, out); diff != "" {
			t.Fatalf("Error matching output and expected: \n%s", diff)
		}
	}
}
