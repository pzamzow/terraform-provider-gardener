package shoot

import (
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
)

func NetworkingResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Description: "Type identifies the type of the networking plugin.",
				Required:    true,
			},
			"pods": {
				Type:        schema.TypeString,
				Description: "Pods is the CIDR of the pod network.",
				Optional:    true,
			},
			"nodes": {
				Type:        schema.TypeString,
				Description: "Nodes is the CIDR of the entire node network.",
				Optional:    true,
			},
			"services": {
				Type:        schema.TypeString,
				Description: "Services is the CIDR of the service network.",
				Optional:    true,
			},
		},
	}
}

func ExpandNetworking(networking []interface{}) *corev1beta1.Networking {
	obj := &corev1beta1.Networking{}
	if len(networking) == 0 || networking[0] == nil {
		return obj
	}

	in := networking[0].(map[string]interface{})
	if v, ok := in["type"].(string); ok && len(v) > 0 {
		obj.Type = v
	}

	if v, ok := in["pods"].(string); ok && len(v) > 0 {
		obj.Pods = &v
	}

	if v, ok := in["nodes"].(string); ok && len(v) > 0 {
		obj.Nodes = &v
	}

	if v, ok := in["services"].(string); ok && len(v) > 0 {
		obj.Services = &v
	}

	return obj
}

func FlattenNetworking(in *corev1beta1.Networking) []interface{} {
	att := make(map[string]interface{})

	if in.Nodes != nil {
		att["nodes"] = *in.Nodes
	}
	if in.Pods != nil {
		att["pods"] = *in.Pods
	}
	if in.Services != nil {
		att["services"] = *in.Services
	}
	if len(in.Type) > 0 {
		att["type"] = in.Type
	}

	return []interface{}{att}
}
