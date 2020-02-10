package gardener

import (
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
)

func ProviderResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Description: "Type is the type of the provider.",
				Required:    true,
			},
			"infrastructure_config": {
				Type:        schema.TypeList,
				Description: "InfrastructureConfig contains the provider-specific infrastructure config blob.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"azure": {
							Type:        schema.TypeList,
							Description: "Azure contains the Shoot specification for the Azure Cloud Platform.",
							Optional:    true,
							MaxItems:    1,
							Elem:        AzureResource(),
						},
					},
				},
			},
			"worker": {
				Type:             schema.TypeList,
				Description:      "Workers is a list of worker groups.",
				Required:         true,
				Elem:             WorkerResource(),
				DiffSuppressFunc: suppressMissingOptionalConfigurationBlock,
			},
		},
	}
}

func ExpandProvider(provider []interface{}) corev1beta1.Provider {
	obj := corev1beta1.Provider{}
	if len(provider) == 0 || provider[0] == nil {
		return obj
	}
	in := provider[0].(map[string]interface{})

	if v, ok := in["type"].(string); ok && len(v) > 0 {
		obj.Type = v
	}

	if v, ok := in["infrastructure_config"].([]interface{}); ok && len(v) > 0 {
		cloud := v[0].(map[string]interface{})
		if az, ok := cloud["azure"].([]interface{}); ok && len(az) > 0 {
			obj.InfrastructureConfig = ExpandAzure(az)
		}
	}
	if workers, ok := in["worker"].([]interface{}); ok && len(workers) > 0 {
		for _, w := range workers {
			if w, ok := w.(map[string]interface{}); ok {
				workerObj := ExpandWorker(w)
				obj.Workers = append(obj.Workers, workerObj)
			}
		}
	}

	return obj
}

func FlattenProvider(in corev1beta1.Provider) []interface{} {
	att := make(map[string]interface{})

	if len(in.Type) > 0 {
		att["type"] = in.Type
	}

	if len(in.Workers) > 0 {

	}

	if in.InfrastructureConfig != nil {
		att["infrastructure_config"] = flattenInfrastructureConfig(in.Type, in.InfrastructureConfig)
	}

	return []interface{}{att}
}

func flattenInfrastructureConfig(providerType string, in *corev1beta1.ProviderConfig) []interface{} {
	att := map[string]interface{}{}

	if providerType == "azure" {
		att["azure"] = FlattenAzure(in)
	}

	return []interface{}{att}
}