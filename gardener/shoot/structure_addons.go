package shoot

import (
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
	lib "github.com/kyma-incubator/terraform-provider-gardener/gardener/lib"
)

func AddonsResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"kubernetes_dashboard": {
				Type:             schema.TypeList,
				Description:      "Kubernetes dashboard holds configuration settings for the kubernetes dashboard addon.",
				Optional:         true,
				MaxItems:         1,
				DiffSuppressFunc: lib.SuppressMissingOptionalConfigurationBlock,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"authentication_mode": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: lib.SuppressEmptyNewValue,
							Default:          "token",
						},
					},
				},
			},
			"nginx_ingress": {
				Type:             schema.TypeList,
				Description:      "NginxIngress holds configuration settings for the nginx-ingress addon.",
				Optional:         true,
				MaxItems:         1,
				DiffSuppressFunc: lib.SuppressMissingOptionalConfigurationBlock,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
		},
	}
}

func ExpandAddons(addon []interface{}) *corev1beta1.Addons {
	obj := &corev1beta1.Addons{}

	if len(addon) == 0 || addon[0] == nil {
		return obj
	}
	in := addon[0].(map[string]interface{})

	if v, ok := in["kubernetes_dashboard"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.KubernetesDashboard = &corev1beta1.KubernetesDashboard{}

		if v, ok := v["enabled"].(bool); ok {
			obj.KubernetesDashboard.Enabled = v
		}
		if v, ok := v["authentication_mode"].(string); ok && len(v) > 0 {
			obj.KubernetesDashboard.AuthenticationMode = &v
		}
	}
	if v, ok := in["nginx_ingress"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.NginxIngress = &corev1beta1.NginxIngress{}

		if v, ok := v["enabled"].(bool); ok {
			obj.NginxIngress.Enabled = v
		}
	}

	return obj
}

func FlattenAddons(in *corev1beta1.Addons) []interface{} {
	att := make(map[string]interface{})

	if in.KubernetesDashboard != nil {
		dashboard := make(map[string]interface{})
		dashboard["enabled"] = in.KubernetesDashboard.Enabled
		if in.KubernetesDashboard.AuthenticationMode != nil {
			dashboard["authentication_mode"] = *in.KubernetesDashboard.AuthenticationMode
		}
		att["kubernetes_dashboard"] = []interface{}{dashboard}
	}
	if in.NginxIngress != nil {
		ingress := make(map[string]interface{})
		ingress["enabled"] = in.NginxIngress.Enabled
		att["nginx_ingress"] = []interface{}{ingress}
	}

	return []interface{}{att}
}
