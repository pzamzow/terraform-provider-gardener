package gardener

import (
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
)

func MonitoringResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"alerting": {
				Type:             schema.TypeList,
				Description:      "Alert configuration to send notification to email lists.",
				Optional:         true,
				MaxItems:         1,
				DiffSuppressFunc: suppressMissingOptionalConfigurationBlock,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"emailreceivers": {
							Type:             schema.TypeSet,
							Description:      "List of people who receiving alerts for this shoots",
							Optional:         true,
							Elem:             &schema.Schema{Type: schema.TypeString},
							Set:              schema.HashString,
							DiffSuppressFunc: suppressEmptyNewValue,
						},
					},
				},
			},
		},
	}
}

func ExpandMonitoring(monitoring []interface{}) *corev1beta1.Monitoring {
	obj := corev1beta1.Monitoring{}
	if len(monitoring) == 0 || monitoring[0] == nil {
		return &obj
	}

	in := monitoring[0].(map[string]interface{})
	if v, ok := in["alerting"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
		alert := corev1beta1.Alerting{}

		in := v[0].(map[string]interface{})
		if v, ok := in["emailreceivers"].(*schema.Set); ok {
			alert.EmailReceivers = expandSet(v)
		}

		obj.Alerting = &alert
	}

	return &obj
}

func FlattenMonitoring(in *corev1beta1.Monitoring) []interface{} {
	att := make(map[string]interface{})

	if in.Alerting != nil {
		alerting := make(map[string]interface{})

		if len(in.Alerting.EmailReceivers) > 0 {
			alerting["emailreceivers"] = in.Alerting.EmailReceivers
		}

		att["alerting"] = []interface{}{alerting}
	}

	return []interface{}{att}
}
