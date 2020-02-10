package gardener

import (
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
)

func KubeletResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"feature_gates": {
				Type:        schema.TypeMap,
				Description: "FeatureGates contains information about enabled feature gates.",
				Optional:    true,
			},
			"pod_pids_limit": {
				Type:        schema.TypeInt,
				Description: "PodPIDsLimit is the maximum number of process IDs per pod allowed by the kubelet.",
				Optional:    true,
			},
			"cpu_cfs_quota": {
				Type:        schema.TypeBool,
				Description: "CPUCFSQuota allows you to disable/enable CPU throttling for Pods.",
				Optional:    true,
			},
			"cpu_manager_policy": {
				Type:        schema.TypeString,
				Description: "CPUManagerPolicy allows to set alternative CPU management policies (default: none).",
				Optional:    true,
			},
		},
	}
}

func ExpandKubelet(kubelet []interface{}) *corev1beta1.KubeletConfig {
	obj := &corev1beta1.KubeletConfig{}

	if len(kubelet) == 0 || kubelet[0] == nil {
		return obj
	}
	in := kubelet[0].(map[string]interface{})

	if v, ok := in["feature_gates"].(map[string]interface{}); ok {
		obj.FeatureGates = expandBoolMap(v)
	}
	if v, ok := in["pod_pids_limit"].(*int64); ok {
		obj.PodPIDsLimit = v
	}
	if v, ok := in["cpu_cfs_quota"].(*bool); ok {
		obj.CPUCFSQuota = v
	}

	return obj
}

func FlattenKubelet(in *corev1beta1.KubeletConfig) []interface{} {
	att := make(map[string]interface{})

	if in.FeatureGates != nil {
		att["feature_gates"] = in.FeatureGates
	}
	if in.PodPIDsLimit != nil {
		att["pod_pids_limit"] = in.PodPIDsLimit
	}
	if in.CPUCFSQuota != nil {
		att["cpu_cfs_quota"] = in.CPUCFSQuota
	}

	return []interface{}{att}
}
