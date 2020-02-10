package gardener

import (
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
)

func ClusterAutoscalerResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"scale_down_utilization_threshold": {
				Type:        schema.TypeFloat,
				Description: "ScaleDownUtilizationThreshold defines the threshold in % under which a node is being removed.",
				Optional:    true,
			},
			"scale_down_unneeded_time": {
				Type:        schema.TypeString,
				Description: "ScaleDownUnneededTime defines how long a node should be unneeded before it is eligible for scale down (default: 10 mins).",
				Optional:    true,
			},
			"scale_down_delay_after_add": {
				Type:        schema.TypeString,
				Description: "ScaleDownDelayAfterAdd defines how long after scale up that scale down evaluation resumes (default: 10 mins).",
				Optional:    true,
			},
			"scale_down_delay_after_failure": {
				Type:        schema.TypeString,
				Description: "ScaleDownDelayAfterFailure how long after scale down failure that scale down evaluation resumes (default: 3 mins).",
				Optional:    true,
			},
			"scale_down_delay_after_delete": {
				Type:        schema.TypeString,
				Description: "ScaleDownDelayAfterDelete how long after node deletion that scale down evaluation resumes, defaults to scanInterval (defaults to ScanInterval).",
				Optional:    true,
			},
			"scan_interval": {
				Type:        schema.TypeString,
				Description: "ScanInterval how often cluster is reevaluated for scale up or down (default: 10 secs).",
				Optional:    true,
			},
		},
	}
}

func ExpandClusterAutoscaler(autoscaler []interface{}) *corev1beta1.ClusterAutoscaler {
	obj := corev1beta1.ClusterAutoscaler{}

	if len(autoscaler) == 0 || autoscaler[0] == nil {
		return &obj
	}
	in := autoscaler[0].(map[string]interface{})

	if v, ok := in["scale_down_utilization_threshold"].(float64); ok {
		obj.ScaleDownUtilizationThreshold = &v
	}

	return &obj
}

func FlattenClusterAutoscaler(in *corev1beta1.ClusterAutoscaler) []interface{} {
	att := make(map[string]interface{})

	if in.ScaleDownUtilizationThreshold != nil {
		att["scale_down_utilization_threshold"] = in.ScaleDownUtilizationThreshold
	}

	return []interface{}{att}
}
