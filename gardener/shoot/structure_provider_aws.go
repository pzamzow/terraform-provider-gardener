package shoot

import (
	"encoding/json"

	awsAlpha1 "github.com/gardener/gardener-extension-provider-aws/pkg/apis/aws/v1alpha1"
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
	lib "github.com/kyma-incubator/terraform-provider-gardener/gardener/lib"
)

func AWSInfrastructureConfigResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enableecraccess": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"networks": {
				Type:        schema.TypeList,
				Description: "Networks is the network configuration (VNet, subnets, etc.).",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc": {
							Type:             schema.TypeList,
							Description:      "VPC ID or CIDR for aws",
							Required:         true,
							MaxItems:         1,
							DiffSuppressFunc: lib.SuppressMissingOptionalConfigurationBlock,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Description: "ID of the VPC.",
										Optional:    true,
									},
									"cidr": {
										Type:        schema.TypeString,
										Description: "CIDR is the VPC CIDR.",
										Optional:    true,
									},
								},
							},
						},
						"zones": {
							Type:             schema.TypeList,
							Description:      "List of zones.",
							Optional:         true,
							DiffSuppressFunc: lib.SuppressMissingOptionalConfigurationBlock,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Description: "Name is the zone name.",
										Optional:    true,
									},
									"internal": {
										Type:        schema.TypeString,
										Description: "internal CIDR",
										Optional:    true,
									},
									"public": {
										Type:        schema.TypeString,
										Description: "public cidr",
										Optional:    true,
									},
									"workers": {
										Type:        schema.TypeString,
										Description: "worker cidr",
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

func ExpandAWSInfrastructureConfig(az []interface{}) *corev1beta1.ProviderConfig {
	obj := corev1beta1.ProviderConfig{}
	if len(az) == 0 && az[0] == nil {
		return &obj
	}
	in := az[0].(map[string]interface{})

	awsConfigObj := awsAlpha1.InfrastructureConfig{}
	awsConfigObj.APIVersion = "aws.provider.extensions.gardener.cloud/v1alpha1"
	awsConfigObj.Kind = "InfrastructureConfig"
	if v, ok := in["enableecraccess"].(bool); ok {
		awsConfigObj.EnableECRAccess = &v
	}
	if v, ok := in["networks"].([]interface{}); ok && len(v) > 0 {
		awsConfigObj.Networks = expandAWSNetworks(v)
	}
	obj.Raw, _ = json.Marshal(awsConfigObj)
	return &obj
}

func expandAWSNetworks(networks []interface{}) awsAlpha1.Networks {
	obj := awsAlpha1.Networks{}
	if networks == nil {
		return obj
	}
	in := networks[0].(map[string]interface{})

	if v, ok := in["vpc"].([]interface{}); ok {
		obj.VPC = expandAWSVPC(v)
	}

	if v, ok := in["zones"].([]interface{}); ok {
		obj.Zones = expandAwsZones(v)
	}

	return obj
}

func expandAWSVPC(vpc []interface{}) awsAlpha1.VPC {
	obj := awsAlpha1.VPC{}

	if len(vpc) == 0 && vpc[0] == nil {
		return obj
	}
	in := vpc[0].(map[string]interface{})

	if v, ok := in["id"].(string); ok && len(v) > 0 {
		obj.ID = &v
	}
	if v, ok := in["cidr"].(string); ok && len(v) > 0 {
		obj.CIDR = &v
	}
	return obj
}

func expandAwsZones(zones []interface{}) []awsAlpha1.Zone {
	result := make([]awsAlpha1.Zone, len(zones))
	for i, k := range zones {
		z := awsAlpha1.Zone{}
		if v, ok := k.(map[string]interface{})["name"].(string); ok && len(v) > 0 {
			z.Name = v
		}
		if v, ok := k.(map[string]interface{})["internal"].(string); ok && len(v) > 0 {
			z.Internal = v
		}
		if v, ok := k.(map[string]interface{})["public"].(string); ok && len(v) > 0 {
			z.Public = v
		}
		if v, ok := k.(map[string]interface{})["workers"].(string); ok && len(v) > 0 {
			z.Workers = v
		}

		result[i] = z
	}
	return result
}

func FlattenAWSInfrastructureConfig(in *corev1beta1.ProviderConfig) []interface{} {
	att := make(map[string]interface{})
	net := make(map[string]interface{})
	vpc := make(map[string]interface{})

	awsConfigObj := awsAlpha1.InfrastructureConfig{}
	if err := json.Unmarshal(in.RawExtension.Raw, &awsConfigObj); err != nil {
		return []interface{}{att}
	}

	if awsConfigObj.EnableECRAccess != nil {
		att["enableecraccess"] = *awsConfigObj.EnableECRAccess
	}
	if awsConfigObj.Networks.VPC.ID != nil {
		vpc["id"] = *awsConfigObj.Networks.VPC.ID
	}
	if awsConfigObj.Networks.VPC.CIDR != nil {
		vpc["cidr"] = *awsConfigObj.Networks.VPC.CIDR
	}
	net["vpc"] = []interface{}{vpc}

	if len(awsConfigObj.Networks.Zones) > 0 {
		zones := make([]interface{}, len(awsConfigObj.Networks.Zones))
		for i, v := range awsConfigObj.Networks.Zones {
			zone := map[string]interface{}{}
			if len(v.Name) > 0 {
				zone["name"] = v.Name
			}
			if len(v.Internal) > 0 {
				zone["internal"] = v.Internal
			}
			if len(v.Public) > 0 {
				zone["public"] = v.Public
			}
			if len(v.Workers) > 0 {
				zone["workers"] = v.Workers
			}
			zones[i] = zone
		}
		net["zones"] = zones
	}
	att["networks"] = []interface{}{net}

	return []interface{}{att}
}
