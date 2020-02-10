package gardener

import (
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
)

func MaintenanceResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"auto_update": {
				Type:             schema.TypeList,
				Description:      "AutoUpdate contains information about which constraints should be automatically updated.",
				MaxItems:         1,
				Optional:         true,
				DiffSuppressFunc: suppressMissingOptionalConfigurationBlock,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kubernetes_version": {
							Type:        schema.TypeBool,
							Description: "KubernetesVersion indicates whether the patch Kubernetes version may be automatically updated.",
							Optional:    true,
							Default:     true,
						},
						"machine_image_version": {
							Type:        schema.TypeBool,
							Description: "machineImageVersion indicates whether the machine image version may be automatically updated.",
							Optional:    true,
							Default:     true,
						},
					},
				},
			},
			"time_window": {
				Type:             schema.TypeList,
				Description:      "TimeWindow contains information about the time window for maintenance operations.",
				Optional:         true,
				MaxItems:         1,
				DiffSuppressFunc: suppressMissingOptionalConfigurationBlock,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"begin": {
							Type:             schema.TypeString,
							Description:      "Begin is the beginning of the time window in the format HHMMSS+ZONE, e.g. '220000+0100'.",
							Optional:         true,
							DiffSuppressFunc: suppressEmptyNewValue,
						},
						"end": {
							Type:             schema.TypeString,
							Description:      "End is the end of the time window in the format HHMMSS+ZONE, e.g. '220000+0100'.",
							Optional:         true,
							DiffSuppressFunc: suppressEmptyNewValue,
						},
					},
				},
			},
		},
	}
}

func ExpandMaintenance(maintenance []interface{}) *corev1beta1.Maintenance {
	obj := &corev1beta1.Maintenance{}

	if len(maintenance) == 0 || maintenance[0] == nil {
		return obj
	}
	in := maintenance[0].(map[string]interface{})

	if v, ok := in["auto_update"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.AutoUpdate = &corev1beta1.MaintenanceAutoUpdate{}

		if v, ok := v["kubernetes_version"].(bool); ok {
			obj.AutoUpdate.KubernetesVersion = v
		}
		if v, ok := v["machine_image_version"].(bool); ok {
			obj.AutoUpdate.MachineImageVersion = v
		}
	}
	if v, ok := in["time_window"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.TimeWindow = &corev1beta1.MaintenanceTimeWindow{}

		if v, ok := v["begin"].(string); ok && len(v) > 0 {
			obj.TimeWindow.Begin = v
		}
		if v, ok := v["end"].(string); ok && len(v) > 0 {
			obj.TimeWindow.End = v
		}
	}

	return obj
}

func FlattenMaintenance(in *corev1beta1.Maintenance) []interface{} {
	att := make(map[string]interface{})

	if in.AutoUpdate != nil {
		update := make(map[string]interface{})
		update["kubernetes_version"] = in.AutoUpdate.KubernetesVersion
		update["machine_image_version"] = in.AutoUpdate.MachineImageVersion
		att["auto_update"] = []interface{}{update}
	}
	if in.TimeWindow != nil {
		window := make(map[string]interface{})
		if len(in.TimeWindow.Begin) > 0 {
			window["begin"] = in.TimeWindow.Begin
		}
		if len(in.TimeWindow.End) > 0 {
			window["end"] = in.TimeWindow.End
		}
		att["time_window"] = []interface{}{window}
	}

	return []interface{}{att}
}
