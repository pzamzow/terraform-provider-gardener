package gardener

import (
	"encoding/json"

	azAlpha1 "github.com/gardener/gardener-extensions/controllers/provider-azure/pkg/apis/azure/v1alpha1"
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
)

func AzureResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"networks": {
				Type:        schema.TypeList,
				Description: "NetworkConfig holds information about the Kubernetes and infrastructure networks.",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"workers": {
							Type:        schema.TypeString,
							Description: "Workers is the worker subnet range to create (used for the VMs).",
							Required:    true,
						},
						"service_endpoints": {
							Type:        schema.TypeSet,
							Description: "List of Azure service endpoints connect to the created VNet.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         schema.HashString,
						},
						"vnet": {
							Type:        schema.TypeList,
							Description: "VNet indicates whether to use an existing VNet or create a new one.",
							Required:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Description: "Name is the VNet name.",
										Optional:    true,
									},
									"cidr": {
										Type:        schema.TypeString,
										Description: "CIDR is the VNet CIDR.",
										Optional:    true,
									},
									"resource_group": {
										Type:        schema.TypeString,
										Description: "ResourceGroup is the resource group where the existing vNet belongs to.",
										Optional:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func ExpandAzure(az []interface{}) *corev1beta1.ProviderConfig {
	obj := corev1beta1.ProviderConfig{}
	if len(az) == 0 && az[0] == nil {
		return &obj
	}
	in := az[0].(map[string]interface{})

	azConfigObj := azAlpha1.InfrastructureConfig{}
	azConfigObj.APIVersion = "azure.provider.extensions.gardener.cloud/v1alpha1"
	azConfigObj.Kind = "InfrastructureConfig"
	if v, ok := in["networks"].([]interface{}); ok && len(v) > 0 {
		azConfigObj.Networks = expandNetworks(v)
	}
	obj.Raw, _ = json.Marshal(azConfigObj)
	return &obj
}

func expandNetworks(networks []interface{}) azAlpha1.NetworkConfig {
	obj := azAlpha1.NetworkConfig{}
	if networks == nil {
		return obj
	}
	in := networks[0].(map[string]interface{})

	if v, ok := in["vnet"].([]interface{}); ok {
		obj.VNet = expandVNET(v)
	}
	if v, ok := in["workers"].(string); ok {
		obj.Workers = v
	}
	if v, ok := in["service_endpoints"].(*schema.Set); ok {
		obj.ServiceEndpoints = expandSet(v)
	}

	return obj
}

func expandVNET(vnet []interface{}) azAlpha1.VNet {
	obj := azAlpha1.VNet{}

	if len(vnet) == 0 && vnet[0] == nil {
		return obj
	}
	in := vnet[0].(map[string]interface{})

	if v, ok := in["name"].(string); ok && len(v) > 0 {
		obj.Name = &v
	}
	if v, ok := in["resource_group"].(string); ok && len(v) > 0 {
		obj.ResourceGroup = &v
	}

	if v, ok := in["cidr"].(string); ok && len(v) > 0 {
		obj.CIDR = &v
	}
	return obj
}

func FlattenAzure(in *corev1beta1.ProviderConfig) []interface{} {
	att := make(map[string]interface{})

	azConfigObj := azAlpha1.InfrastructureConfig{}
	if err := json.Unmarshal(in.RawExtension.Raw, &azConfigObj); err != nil {
		return []interface{}{att}
	}

	net := make(map[string]interface{})
	if len(azConfigObj.Networks.Workers) > 0 {
		net["workers"] = azConfigObj.Networks.Workers
	}
	if len(azConfigObj.Networks.ServiceEndpoints) > 0 {
		net["service_endpoints"] = azConfigObj.Networks.ServiceEndpoints
	}
	vnet := make(map[string]interface{})
	if azConfigObj.Networks.VNet.CIDR != nil {
		vnet["cidr"] = *azConfigObj.Networks.VNet.CIDR
	}
	if azConfigObj.Networks.VNet.Name != nil {
		vnet["name"] = *azConfigObj.Networks.VNet.Name
	}
	if azConfigObj.Networks.VNet.ResourceGroup != nil {
		vnet["resource_group"] = *azConfigObj.Networks.VNet.ResourceGroup
	}
	net["vnet"] = []interface{}{vnet}
	att["networks"] = []interface{}{net}

	return []interface{}{att}
}
