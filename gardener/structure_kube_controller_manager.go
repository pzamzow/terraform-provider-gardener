package gardener

import (
	"time"

	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
)


func KubeControllerManagerResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"node_cidr_mask_size": {
				Type:        schema.TypeInt,
				Description: "Size of the node mask.",
				Optional:    true,
			},
		},
	}
}

func ExpandKubeControllerManager(controller []interface{}) *corev1beta1.KubeControllerManagerConfig {
	obj := &corev1beta1.KubeControllerManagerConfig{}

	if len(controller) == 0 || controller[0] == nil {
		return obj
	}
	in := controller[0].(map[string]interface{})

	if v, ok := in["feature_gates"].(map[string]interface{}); ok {
		obj.FeatureGates = expandBoolMap(v)
	}
	if v, ok := in["horizontal_pod_autoscaler"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.HorizontalPodAutoscalerConfig = &corev1beta1.HorizontalPodAutoscalerConfig{}

		if v, ok := v["downscale_delay"].(string); ok && len(v) > 0 {
			obj.HorizontalPodAutoscalerConfig.DownscaleDelay = expandDuration(v)
		}
		if v, ok := v["sync_period"].(string); ok && len(v) > 0 {
			obj.HorizontalPodAutoscalerConfig.SyncPeriod = expandDuration(v)
		}
		if v, ok := v["tolerance"].(*float64); ok {
			obj.HorizontalPodAutoscalerConfig.Tolerance = v
		}
		if v, ok := v["upscale_delay"].(string); ok && len(v) > 0 {
			obj.HorizontalPodAutoscalerConfig.UpscaleDelay = expandDuration(v)
		}
		if v, ok := v["downscale_stabilization"].(string); ok && len(v) > 0 {
			obj.HorizontalPodAutoscalerConfig.DownscaleStabilization = expandDuration(v)
		}
		if v, ok := v["initial_readiness_delay"].(string); ok && len(v) > 0 {
			obj.HorizontalPodAutoscalerConfig.InitialReadinessDelay = expandDuration(v)
		}
		if v, ok := v["cpu_initialization_period"].(string); ok && len(v) > 0 {
			obj.HorizontalPodAutoscalerConfig.CPUInitializationPeriod = expandDuration(v)
		}
	}

	return obj
}

func expandDuration(v string) *corev1beta1.GardenerDuration {
	d, err := time.ParseDuration(v)
	if err != nil {
		return &corev1beta1.GardenerDuration{
			Duration: d,
		}
	}

	return nil
}

func FlattenKubeControllerManager(in *corev1beta1.KubeControllerManagerConfig) []interface{} {
	att := make(map[string]interface{})

	if in.FeatureGates != nil {
		att["node_cidr_mask_size"] = in.NodeCIDRMaskSize
	}

	return []interface{}{att}
}
