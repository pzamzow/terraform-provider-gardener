package gardener

import (
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
)

func VolumeResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Description: "VolumeType is the type of the root volumes.",
				Required:    true,
			},
			"size": {
				Type:        schema.TypeString,
				Description: "VolumeSize is the size of the root volume.",
				Required:    true,
			},
		},
	}
}

func ExpandVolume(v []interface{}) *corev1beta1.Volume {
	obj := &corev1beta1.Volume{}

	if len(v) == 0 && v[0] == nil {
		return obj
	}
	in := v[0].(map[string]interface{})

	if c, ok := in["type"].(string); ok && len(c) > 0 {
		obj.Type = &c
	}

	if c, ok := in["size"].(string); ok && len(c) > 0 {
		obj.Size = c
	}
	return obj
}

func FlattenVolume(in *corev1beta1.Volume) []interface{} {
	att := map[string]interface{}{}

	if len(in.Size) > 0 {
		att["size"] = in.Size
	}
	if in.Type != nil {
		att["type"] = *in.Type
	}

	return []interface{}{att}
}
