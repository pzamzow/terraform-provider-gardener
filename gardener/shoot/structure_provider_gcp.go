package shoot

import (
	"encoding/json"

	gcpAlpha1 "github.com/gardener/gardener-extension-provider-gcp/pkg/apis/gcp/v1alpha1"
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
	lib "github.com/kyma-incubator/terraform-provider-gardener/gardener/lib"
)

func GCPControlPlaneConfigResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"zone": {
				Type:        schema.TypeString,
				Description: "Zone is the GCP zone.",
				Required:    true,
			},
		},
	}
}

func GCPInfrastructureConfigResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"networks": {
				Type:        schema.TypeList,
				Description: "Networks is the network configuration (VPC, subnets, etc.)",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc": {
							Type:             schema.TypeList,
							Description:      "VPC indicates whether to use an existing VPC or create a new one.",
							Optional:         true,
							MaxItems:         1,
							DiffSuppressFunc: lib.SuppressMissingOptionalConfigurationBlock,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Description: "Name is the VPC name.",
										Optional:    true,
									},
									"cloud_router": {
										Type:             schema.TypeList,
										Description:      "CloudRouter indicates whether to use an existing CloudRouter or create a new one.",
										Optional:         true,
										MaxItems:         1,
										DiffSuppressFunc: lib.SuppressMissingOptionalConfigurationBlock,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Description: "Name is the CloudRouter name.",
													Optional:    true,
												},
											},
										},
									},
								},
							},
						},
						"cloud_nat": {
							Type:             schema.TypeList,
							Description:      "CloudNAT contains configuration about the the CloudNAT configuration",
							Optional:         true,
							MaxItems:         1,
							DiffSuppressFunc: lib.SuppressMissingOptionalConfigurationBlock,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"min_ports_per_vm": {
										Type:        schema.TypeInt,
										Description: "MinPortsPerVM is the minimum number of ports allocated to a VM in the NAT config. The default value is 2048 ports.",
										Optional:    true,
									},
								},
							},
						},
						"internal": {
							Type:        schema.TypeString,
							Description: "Internal is a private subnet (used for internal load balancers).",
							Optional:    true,
						},
						"workers": {
							Type:        schema.TypeString,
							Description: "Workers is the worker subnet range to create (used for the VMs).",
							Required:    true,
						},
						"flow_logs": {
							Type:             schema.TypeList,
							Description:      "FlowLogs contains the flow log configuration for the subnet.",
							Optional:         true,
							MaxItems:         1,
							DiffSuppressFunc: lib.SuppressMissingOptionalConfigurationBlock,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"aggregation_interval": {
										Type:        schema.TypeString,
										Description: "AggregationInterval for collecting flow logs.",
										Optional:    true,
									},
									"flow_sampling": {
										Type:        schema.TypeFloat,
										Description: "FlowSampling sets the sampling rate of VPC flow logs within the subnetwork where 1.0 means all collected logs are reported and 0.0 means no logs are reported.",
										Optional:    true,
									},
									"metadata": {
										Type:        schema.TypeString,
										Description: "Metadata configures whether metadata fields should be added to the reported VPC flow logs.",
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

func ExpandGCPControlPlaneConfig(gcp []interface{}) *corev1beta1.ProviderConfig {
	gcpConfigObj := gcpAlpha1.ControlPlaneConfig{}
	obj := corev1beta1.ProviderConfig{}
	if len(gcp) == 0 && gcp[0] == nil {
		return &obj
	}
	in := gcp[0].(map[string]interface{})

	gcpConfigObj.APIVersion = "gcp.provider.extensions.gardener.cloud/v1alpha1"
	gcpConfigObj.Kind = "ControlPlaneConfig"

	if v, ok := in["zone"].(string); ok && len(v) > 0 {
		gcpConfigObj.Zone = v
	}

	obj.Raw, _ = json.Marshal(gcpConfigObj)
	return &obj
}

func ExpandGCPInfrastructureConfig(az []interface{}) *corev1beta1.ProviderConfig {
	obj := corev1beta1.ProviderConfig{}
	if len(az) == 0 && az[0] == nil {
		return &obj
	}
	in := az[0].(map[string]interface{})

	gcpConfigObj := gcpAlpha1.InfrastructureConfig{}
	gcpConfigObj.APIVersion = "gcp.provider.extensions.gardener.cloud/v1alpha1"
	gcpConfigObj.Kind = "InfrastructureConfig"
	if v, ok := in["networks"].([]interface{}); ok && len(v) > 0 {
		gcpConfigObj.Networks = expandGCPNetworks(v)
	}
	obj.Raw, _ = json.Marshal(gcpConfigObj)
	return &obj
}

func expandGCPNetworks(networks []interface{}) gcpAlpha1.NetworkConfig {
	obj := gcpAlpha1.NetworkConfig{}
	if networks == nil {
		return obj
	}
	in := networks[0].(map[string]interface{})

	if v, ok := in["vpc"].([]interface{}); ok && len(v) > 0 {
		obj.VPC = expandGCPVPC(v)
	}
	if v, ok := in["workers"].(string); ok && len(v) > 0 {
		obj.Workers = v
	}
	if v, ok := in["internal"].(string); ok && len(v) > 0 {
		obj.Internal = &v
	}
	if v, ok := in["cloud_nat"].([]interface{}); ok && len(v) > 0 {
		obj.CloudNAT = expandGCPCloudNat(v)
	}
	if v, ok := in["flow_logs"].([]interface{}); ok && len(v) > 0 {
		obj.FlowLogs = expandGCPFlowLogs(v)
	}

	return obj
}

func expandGCPVPC(vpc []interface{}) *gcpAlpha1.VPC {
	obj := gcpAlpha1.VPC{}
	if len(vpc) == 0 && vpc[0] == nil {
		return &obj
	}
	in := vpc[0].(map[string]interface{})

	if v, ok := in["name"].(string); ok && len(v) > 0 {
		obj.Name = v
	}
	if v, ok := in["cloud_router"].([]interface{}); ok && len(v) > 0 {
		obj.CloudRouter = expandGCPCloudRouter(v)
	}
	return &obj
}

func expandGCPCloudRouter(cr []interface{}) *gcpAlpha1.CloudRouter {
	obj := gcpAlpha1.CloudRouter{}
	if len(cr) == 0 && cr[0] == nil {
		return &obj
	}
	in := cr[0].(map[string]interface{})

	if v, ok := in["name"].(string); ok && len(v) > 0 {
		obj.Name = v
	}
	return &obj
}

func expandGCPCloudNat(cn []interface{}) *gcpAlpha1.CloudNAT {
	obj := gcpAlpha1.CloudNAT{}
	if len(cn) == 0 && cn[0] == nil {
		return &obj
	}

	in := cn[0].(map[string]interface{})

	if v, ok := in["min_ports_per_vm"].(int); ok {
		f := int32(v)
		obj.MinPortsPerVM = &f
	}
	return &obj
}

func expandGCPFlowLogs(fl []interface{}) *gcpAlpha1.FlowLogs {
	obj := gcpAlpha1.FlowLogs{}
	if len(fl) == 0 && fl[0] == nil {
		return &obj
	}
	in := fl[0].(map[string]interface{})

	if v, ok := in["aggregation_interval"].(string); ok && len(v) > 0 {
		obj.AggregationInterval = &v
	}
	if v, ok := in["flow_sampling"].(float64); ok {
		flowSampling := float32(v)
		obj.FlowSampling = &flowSampling
	}
	if v, ok := in["metadata"].(string); ok {
		obj.Metadata = &v
	}
	return &obj
}

func FlattenGCPControlPlaneConfig(in *corev1beta1.ProviderConfig) []interface{} {
	att := make(map[string]interface{})

	gcpConfigObj := gcpAlpha1.ControlPlaneConfig{}
	if err := json.Unmarshal(in.RawExtension.Raw, &gcpConfigObj); err != nil {
		return []interface{}{att}
	}
	if len(gcpConfigObj.Zone) > 0 {
		att["zone"] = gcpConfigObj.Zone
	}

	return []interface{}{att}
}

func FlattenGCPInfrastructureConfig(in *corev1beta1.ProviderConfig) []interface{} {
	att := make(map[string]interface{})
	net := make(map[string]interface{})

	gcpConfigObj := gcpAlpha1.InfrastructureConfig{}
	if err := json.Unmarshal(in.RawExtension.Raw, &gcpConfigObj); err != nil {
		return []interface{}{att}
	}

	if len(gcpConfigObj.Networks.Workers) > 0 {
		net["workers"] = gcpConfigObj.Networks.Workers
	}

	if gcpConfigObj.Networks.Internal != nil {
		net["internal"] = *gcpConfigObj.Networks.Internal
	}

	vpc := make(map[string]interface{})

	if gcpConfigObj.Networks.VPC != nil && len(gcpConfigObj.Networks.VPC.Name) > 0 {
		vpc["name"] = gcpConfigObj.Networks.VPC.Name
	}
	cr := make(map[string]interface{})
	if gcpConfigObj.Networks.VPC != nil && len(gcpConfigObj.Networks.VPC.CloudRouter.Name) > 0 {
		cr["name"] = gcpConfigObj.Networks.VPC.CloudRouter.Name
		vpc["cloud_router"] = []interface{}{cr}
	}
	net["vpc"] = []interface{}{vpc}

	cn := make(map[string]interface{})
	if gcpConfigObj.Networks.CloudNAT != nil && gcpConfigObj.Networks.CloudNAT.MinPortsPerVM != nil {
		cn["min_ports_per_vm"] = int(*gcpConfigObj.Networks.CloudNAT.MinPortsPerVM)
	}
	net["cloud_nat"] = []interface{}{cn}

	net["flow_logs"] = []interface{}{}
	if gcpConfigObj.Networks.FlowLogs != nil {
		fl := make(map[string]interface{})
		if gcpConfigObj.Networks.FlowLogs.AggregationInterval != nil {
			fl["aggregation_interval"] = *gcpConfigObj.Networks.FlowLogs.AggregationInterval
		}
		if gcpConfigObj.Networks.FlowLogs.Metadata != nil {
			fl["metadata"] = *gcpConfigObj.Networks.FlowLogs.Metadata
		}
		if gcpConfigObj.Networks.FlowLogs.FlowSampling != nil {
			fl["flow_sampling"] = float64(*gcpConfigObj.Networks.FlowLogs.FlowSampling)
		}
		net["flow_logs"] = []interface{}{fl}
	}

	att["networks"] = []interface{}{net}

	return []interface{}{att}
}
