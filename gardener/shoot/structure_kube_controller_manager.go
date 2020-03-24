package shoot

import (
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kyma-incubator/terraform-provider-gardener/gardener/lib"
)

func KubeControllerManagerResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"node_cidr_mask_size": {
				Type:             schema.TypeInt,
				Description:      "Size of the node mask.",
				Optional:         true,
				DiffSuppressFunc: lib.SuppressZeroNewValue,
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

	if v, ok := in["node_cidr_mask_size"].(int); ok {
		size := int32(v)
		obj.NodeCIDRMaskSize = &size
	}

	return obj
}

func FlattenKubeControllerManager(in *corev1beta1.KubeControllerManagerConfig) []interface{} {
	att := make(map[string]interface{})

	if in.NodeCIDRMaskSize != nil {
		size := int(*in.NodeCIDRMaskSize)
		att["node_cidr_mask_size"] = size
	}

	return []interface{}{att}
}
