package gardener

import (
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
)

func HibernationResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Enabled is true if Shoot is hibernated, false otherwise.",
				Required:    true,
			},
			"schedules": {
				Type:        schema.TypeList,
				Description: "Schedules determine the hibernation schedules.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start": {
							Type:        schema.TypeString,
							Description: "Start is a Cron spec at which time a Shoot will be hibernated.",
							Optional:    true,
						},
						"end": {
							Type:        schema.TypeString,
							Description: "End is a Cron spec at which time a Shoot will be woken up.",
							Optional:    true,
						},
						"location": {
							Type:        schema.TypeString,
							Description: "Location is the time location in which both start and and shall be evaluated.",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func ExpandHibernation(hibernation []interface{}) *corev1beta1.Hibernation {
	obj := &corev1beta1.Hibernation{}

	if len(hibernation) == 0 || hibernation[0] == nil {
		return obj
	}
	in := hibernation[0].(map[string]interface{})

	if v, ok := in["enabled"].(bool); ok {
		obj.Enabled = &v
	}
	if schedules, ok := in["schedules"].([]interface{}); ok && len(schedules) > 0 {
		for _, s := range schedules {
			if s, ok := s.(map[string]interface{}); ok {
				scheduleObj := corev1beta1.HibernationSchedule{}

				if v, ok := s["start"].(string); ok && len(v) > 0 {
					scheduleObj.Start = &v
				}
				if v, ok := s["end"].(string); ok && len(v) > 0 {
					scheduleObj.End = &v
				}
				if v, ok := s["location"].(string); ok && len(v) > 0 {
					scheduleObj.Location = &v
				}

				obj.Schedules = append(obj.Schedules, scheduleObj)
			}
		}
	}

	return obj
}

func FlattenHibernation(in *corev1beta1.Hibernation) []interface{} {
	att := make(map[string]interface{})

	if in.Enabled != nil {
		att["enabled"] = *in.Enabled
	}
	if len(in.Schedules) > 0 {
		schedules := make([]interface{}, len(in.Schedules))
		for i, v := range in.Schedules {
			m := map[string]interface{}{}

			if v.Start != nil {
				m["start"] = *v.Start
			}
			if v.End != nil {
				m["end"] = *v.End
			}
			if v.Location != nil {
				m["location"] = *v.Location
			}
			schedules[i] = m
		}
		att["schedules"] = schedules
	}

	return []interface{}{att}
}
