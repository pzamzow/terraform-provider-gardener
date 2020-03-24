package shoot

import (
	"encoding/json"
	"testing"

	azAlpha1 "github.com/gardener/gardener-extension-provider-azure/pkg/apis/azure/v1alpha1"
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	cmp "github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/schema"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func getProviderAzureInfrastructureConfigTestData() (map[string]interface{}, *corev1beta1.ProviderConfig) {
	vnetCIDR := "10.250.0.0/16"
	vnetName := "test"
	resGroup := "test"
	awsConfig, _ := json.Marshal(azAlpha1.InfrastructureConfig{
		TypeMeta: v1.TypeMeta{
			APIVersion: "azure.provider.extensions.gardener.cloud/v1alpha1",
			Kind:       "InfrastructureConfig",
		},
		Zoned: true,
		Networks: azAlpha1.NetworkConfig{
			VNet: azAlpha1.VNet{
				CIDR:          &vnetCIDR,
				Name:          &vnetName,
				ResourceGroup: &resGroup,
			},
			Workers:          "10.250.0.0/19",
			ServiceEndpoints: []string{"microsoft.test"},
		},
	})
	azureRaw := map[string]interface{}{
		"zoned": true,
		"networks": []interface{}{
			map[string]interface{}{
				"vnet": []interface{}{
					map[string]interface{}{
						"cidr":           vnetCIDR,
						"name":           vnetName,
						"resource_group": resGroup,
					},
				},
				"workers":           "10.250.0.0/19",
				"service_endpoints": []interface{}{"microsoft.test"},
			},
		},
	}
	azure := &corev1beta1.ProviderConfig{
		RawExtension: runtime.RawExtension{
			Raw: awsConfig,
		},
	}

	return azureRaw, azure
}

func TestExpandProviderAzure(t *testing.T) {
	azureRaw, azure := getProviderAzureInfrastructureConfigTestData()
	data := schema.TestResourceDataRaw(t, AzureInfrastructureConfigResource().Schema, azureRaw)
	out := ExpandAzureInfrastructureConfig([]interface{}{data.Get("")})
	if diff := cmp.Diff(azure, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}

func TestFlattenProviderAzure(t *testing.T) {
	azureRaw, azure := getProviderAzureInfrastructureConfigTestData()
	out := FlattenAzureInfrastructureConfig(azure)[0].(map[string]interface{})
	data := schema.TestResourceDataRaw(t, AzureInfrastructureConfigResource().Schema, azureRaw).Get("")
	if diff := cmp.Diff(data, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}
