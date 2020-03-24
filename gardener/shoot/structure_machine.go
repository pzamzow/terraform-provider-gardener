package shoot

import (
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
	lib "github.com/kyma-incubator/terraform-provider-gardener/gardener/lib"
)

func MachineResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Description: "Type is the machine type of the worker group.",
				Required:    true,
			},
			"image": {
				Type:             schema.TypeList,
				Description:      "Image holds information about the machine image to use for all nodes of this pool. It will default to the latest version of the first image stated in the referenced CloudProfile if no value has been provided.",
				Optional:         true,
				MaxItems:         1,
				DiffSuppressFunc: lib.SuppressMissingOptionalConfigurationBlock,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Description:      "VolumeSize is the size of the root volume.",
							Optional:         true,
							DiffSuppressFunc: lib.SuppressEmptyNewValue,
						},
						"version": {
							Type:             schema.TypeString,
							Description:      "Version is the version of the shoot's image.",
							Optional:         true,
							DiffSuppressFunc: lib.SuppressEmptyNewValue,
						},
					},
				},
			},
		},
	}
}

func ExpandMachine(m []interface{}) *corev1beta1.Machine {
	obj := &corev1beta1.Machine{}

	if len(m) == 0 && m[0] == nil {
		return obj
	}
	in := m[0].(map[string]interface{})

	if v, ok := in["type"].(string); ok && len(v) > 0 {
		obj.Type = v
	}

	if v, ok := in["image"].([]interface{}); ok && len(v) > 0 {
		obj.Image = expandMachineImage(v)
	}

	return obj
}

func expandMachineImage(si []interface{}) *corev1beta1.ShootMachineImage {
	obj := &corev1beta1.ShootMachineImage{}

	if len(si) == 0 && si[0] == nil {
		return obj
	}
	in := si[0].(map[string]interface{})

	if v, ok := in["name"].(string); ok && len(v) > 0 {
		obj.Name = v
	}

	if v, ok := in["version"].(string); ok && len(v) > 0 {
		obj.Version = &v
	}
	return obj
}

func FlattenMachine(in *corev1beta1.Machine) []interface{} {
	att := map[string]interface{}{}

	if len(in.Type) > 0 {
		att["type"] = in.Type
	}
	if in.Image != nil {
		att["image"] = flattenMachineImage(in.Image)
	}

	return []interface{}{att}
}

func flattenMachineImage(in *corev1beta1.ShootMachineImage) []interface{} {
	att := map[string]interface{}{}

	if len(in.Name) > 0 {
		att["name"] = in.Name
	}
	if in.Version != nil {
		att["version"] = *in.Version
	}

	return []interface{}{att}
}
