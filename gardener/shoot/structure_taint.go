package shoot

import (
	"github.com/hashicorp/terraform/helper/schema"
	v1 "k8s.io/api/core/v1"
)

func TaintResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"effect": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func ExpandTaint(taints []interface{}) []v1.Taint {
	att := make([]v1.Taint, len(taints))

	for i, t := range taints {
		if t, ok := t.(map[string]interface{}); ok {
			taint := v1.Taint{}

			if v, ok := t["key"].(string); ok && len(v) > 0 {
				taint.Key = v
			}
			if v, ok := t["value"].(string); ok && len(v) > 0 {
				taint.Value = v
			}
			if v, ok := t["effect"].(string); ok && len(v) > 0 {
				taint.Effect = v1.TaintEffect(v)
			}

			att[i] = taint
		}
	}

	return att
}

func FlattenTaint(taints []v1.Taint) []interface{} {
	att := make([]interface{}, len(taints))

	for i, v := range taints {
		m := map[string]interface{}{}

		if v.Key != "" {
			m["key"] = v.Key
		}
		if v.Value != "" {
			m["value"] = v.Value
		}
		if v.Effect != "" {
			m["effect"] = string(v.Effect)
		}
		att[i] = m
	}

	return att
}
