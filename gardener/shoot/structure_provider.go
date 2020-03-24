package shoot

import (
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
	lib "github.com/kyma-incubator/terraform-provider-gardener/gardener/lib"
)

func ProviderResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Description: "Type is the type of the provider.",
				Required:    true,
			},
			"control_plane_config": {
				Type:             schema.TypeList,
				Description:      "ControlPlaneConfig contains the provider-specific control plane config blob.",
				Optional:         true,
				MaxItems:         1,
				DiffSuppressFunc: lib.SuppressMissingOptionalConfigurationBlock,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gcp": {
							Type:             schema.TypeList,
							Description:      "GCP contains the Shoot specification for Google Cloud Platform.",
							Optional:         true,
							MaxItems:         1,
							Elem:             GCPControlPlaneConfigResource(),
							DiffSuppressFunc: lib.SuppressMissingOptionalConfigurationBlock,
						},
					},
				},
			},
			"infrastructure_config": {
				Type:             schema.TypeList,
				Description:      "InfrastructureConfig contains the provider-specific infrastructure config blob.",
				Optional:         true,
				MaxItems:         1,
				DiffSuppressFunc: lib.SuppressMissingOptionalConfigurationBlock,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"azure": {
							Type:             schema.TypeList,
							Description:      "Azure contains the Shoot specification for the Azure Cloud Platform.",
							Optional:         true,
							MaxItems:         1,
							Elem:             AzureInfrastructureConfigResource(),
							DiffSuppressFunc: lib.SuppressMissingOptionalConfigurationBlock,
						},
						"aws": {
							Type:             schema.TypeList,
							Description:      "AWS contains the Shoot specification for the AWS Cloud.",
							Optional:         true,
							MaxItems:         1,
							Elem:             AWSInfrastructureConfigResource(),
							DiffSuppressFunc: lib.SuppressMissingOptionalConfigurationBlock,
						},
						"gcp": {
							Type:             schema.TypeList,
							Description:      "AWS contains the Shoot specification for the AWS Cloud.",
							Optional:         true,
							MaxItems:         1,
							Elem:             GCPInfrastructureConfigResource(),
							DiffSuppressFunc: lib.SuppressMissingOptionalConfigurationBlock,
						},
					},
				},
			},
			"worker": {
				Type:             schema.TypeList,
				Description:      "Workers is a list of worker groups.",
				Required:         true,
				Elem:             WorkerResource(),
				DiffSuppressFunc: lib.SuppressMissingOptionalConfigurationBlock,
			},
		},
	}
}

func ExpandProvider(provider []interface{}) *corev1beta1.Provider {
	obj := &corev1beta1.Provider{}
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
			obj.InfrastructureConfig = ExpandAzureInfrastructureConfig(az)
		}
		if aws, ok := cloud["aws"].([]interface{}); ok && len(aws) > 0 {
			obj.InfrastructureConfig = ExpandAWSInfrastructureConfig(aws)
		}
		if aws, ok := cloud["gcp"].([]interface{}); ok && len(aws) > 0 {
			obj.InfrastructureConfig = ExpandGCPInfrastructureConfig(aws)
		}
	}
	if v, ok := in["control_plane_config"].([]interface{}); ok && len(v) > 0 {
		cloud := v[0].(map[string]interface{})
		if gcp, ok := cloud["gcp"].([]interface{}); ok && len(gcp) > 0 {
			obj.ControlPlaneConfig = ExpandGCPControlPlaneConfig(gcp)
		}
	}
	if workers, ok := in["worker"].([]interface{}); ok && len(workers) > 0 {
		for _, w := range workers {
			if w, ok := w.(map[string]interface{}); ok {
				workerObj := ExpandWorker(w)
				obj.Workers = append(obj.Workers, *workerObj)
			}
		}
	}

	return obj
}

func FlattenProvider(in *corev1beta1.Provider) []interface{} {
	att := make(map[string]interface{})

	if len(in.Type) > 0 {
		att["type"] = in.Type
	}

	if len(in.Workers) > 0 {
		workers := make([]interface{}, len(in.Workers))
		for i, w := range in.Workers {
			workers[i] = FlattenWorker(&w)
		}
		att["worker"] = workers
	}

	att["infrastructure_config"] = flattenInfrastructureConfig(in.Type, in.InfrastructureConfig)
	att["control_plane_config"] = flattenControlPlaneConfig(in.Type, in.ControlPlaneConfig)

	return []interface{}{att}
}

func flattenInfrastructureConfig(providerType string, in *corev1beta1.ProviderConfig) []interface{} {
	att := map[string]interface{}{}

	if in == nil {
		return []interface{}{}
	}

	att["azure"] = []interface{}{}
	att["aws"] = []interface{}{}
	att["gcp"] = []interface{}{}

	if providerType == "azure" {
		att["azure"] = FlattenAzureInfrastructureConfig(in)
	} else if providerType == "aws" {
		att["aws"] = FlattenAWSInfrastructureConfig(in)
	} else if providerType == "gcp" {
		att["gcp"] = FlattenGCPInfrastructureConfig(in)
	}

	return []interface{}{att}
}

func flattenControlPlaneConfig(providerType string, in *corev1beta1.ProviderConfig) []interface{} {
	att := map[string]interface{}{}

	if in == nil {
		return []interface{}{}
	}

	att["gcp"] = []interface{}{}

	if providerType == "gcp" {
		att["gcp"] = FlattenGCPControlPlaneConfig(in)
	}

	return []interface{}{att}
}
