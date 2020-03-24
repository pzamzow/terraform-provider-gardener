package shoot

import (
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
	lib "github.com/kyma-incubator/terraform-provider-gardener/gardener/lib"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func WorkerResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"annotations": {
				Type:        schema.TypeMap,
				Description: "Annotations is a map of key/value pairs for annotations for all the `Node` objects in this worker pool.",
				Optional:    true,
			},
			"cabundle": {
				Type:        schema.TypeString,
				Description: "caBundle configuration",
				Optional:    true,
			},
			"kubernetes": {
				Type:             schema.TypeList,
				Description:      "Kubernetes contains configuration for Kubernetes components related to this worker pool.",
				Optional:         true,
				Elem:             WorkerKubernetesResource(),
				DiffSuppressFunc: lib.SuppressMissingOptionalConfigurationBlock,
			},
			"labels": {
				Type:        schema.TypeMap,
				Description: "Labels is a map of key/value pairs for labels for all the `Node` objects in this worker pool.",
				Optional:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name is the name of the worker group.",
				Required:    true,
			},
			"machine": {
				Type:             schema.TypeList,
				Description:      "MachineType is the machine type of the worker group.",
				Required:         true,
				MaxItems:         1,
				Elem:             MachineResource(),
				DiffSuppressFunc: lib.SuppressMissingOptionalConfigurationBlock,
			},
			"minimum": {
				Type:        schema.TypeInt,
				Description: "Minimum is the minimum number of VMs to create.",
				Required:    true,
			},
			"maximum": {
				Type:        schema.TypeInt,
				Description: "Maximum is the maximum number of VMs to create.",
				Required:    true,
			},
			"max_surge": {
				Type:        schema.TypeInt,
				Description: "MaxSurge is maximum number of VMs that are created during an update.",
				Optional:    true,
			},
			"max_unavailable": {
				Type:        schema.TypeInt,
				Description: "MaxUnavailable is the maximum number of VMs that can be unavailable during an update.",
				Optional:    true,
			},
			"taints": {
				Type:             schema.TypeList,
				Description:      "Taints is a list of taints for all the `Node` objects in this worker pool.",
				Optional:         true,
				Elem:             TaintResource(),
				DiffSuppressFunc: lib.SuppressMissingOptionalConfigurationBlock,
			},
			"volume": {
				Type:             schema.TypeList,
				Description:      "Volume contains information about the volume type and size.",
				Required:         true,
				MaxItems:         1,
				Elem:             VolumeResource(),
				DiffSuppressFunc: lib.SuppressMissingOptionalConfigurationBlock,
			},
			"zones": {
				Type:        schema.TypeSet,
				Description: "Zones is a list of availability zones to deploy the Shoot cluster to.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
			},
		},
	}
}

func ExpandWorker(w interface{}) *corev1beta1.Worker {
	obj := &corev1beta1.Worker{}
	if w == nil {
		return obj
	}
	in := w.(map[string]interface{})

	if v, ok := in["annotations"].(map[string]interface{}); ok {
		obj.Annotations = lib.ExpandStringMap(v)
	}

	if v, ok := in["cabundle"].(string); ok && len(v) > 0 {
		obj.CABundle = &v
	}

	if v, ok := in["kubernetes"].([]interface{}); ok && len(v) > 0 {
		obj.Kubernetes = ExpandWorkerKubernetes(v)
	}

	if v, ok := in["labels"].(map[string]interface{}); ok {
		obj.Labels = lib.ExpandStringMap(v)
	}

	if v, ok := in["name"].(string); ok && len(v) > 0 {
		obj.Name = v
	}

	if v, ok := in["machine"].([]interface{}); ok && len(v) > 0 {
		obj.Machine = *ExpandMachine(v)
	}

	if v, ok := in["maximum"].(int); ok {
		obj.Maximum = int32(v)
	}

	if v, ok := in["minimum"].(int); ok {
		obj.Minimum = int32(v)
	}

	if v, ok := in["max_surge"].(int); ok {
		surge := intstr.FromInt(v)
		obj.MaxSurge = &surge
	}

	if v, ok := in["max_unavailable"].(int); ok {
		unavailable := intstr.FromInt(v)
		obj.MaxUnavailable = &unavailable
	}

	if taints, ok := in["taints"].([]interface{}); ok && len(taints) > 0 {
		obj.Taints = ExpandTaint(taints)
	}

	if v, ok := in["volume"].([]interface{}); ok {
		obj.Volume = ExpandVolume(v)
	}

	if v, ok := in["zones"].(*schema.Set); ok {
		obj.Zones = lib.ExpandSet(v)
	}

	return obj
}

func FlattenWorker(in *corev1beta1.Worker) interface{} {
	att := map[string]interface{}{}

	if len(in.Name) > 0 {
		att["name"] = in.Name
	}
	if len(in.Zones) > 0 {
		att["zones"] = lib.NewStringSet(schema.HashString, in.Zones)
	}
	if len(in.Taints) > 0 {
		att["taints"] = FlattenTaint(in.Taints)
	}
	if in.MaxSurge != nil {
		att["max_surge"] = in.MaxSurge.IntValue()
	}
	if in.MaxUnavailable != nil {
		att["max_unavailable"] = in.MaxUnavailable.IntValue()
	}
	if in.CABundle != nil {
		att["cabundle"] = *in.CABundle
	}
	if in.Minimum != 0 {
		att["minimum"] = int(in.Minimum)
	}
	if in.Maximum != 0 {
		att["maximum"] = int(in.Maximum)
	}
	if in.Kubernetes != nil {
		att["kubernetes"] = FlattenWorkerKubernetes(in.Kubernetes)
	}
	if len(in.Annotations) > 0 {
		att["annotations"] = lib.FlattenStringMap(in.Annotations)
	}
	if len(in.Labels) > 0 {
		att["labels"] = lib.FlattenStringMap(in.Labels)
	}
	if in.Volume != nil {
		att["volume"] = FlattenVolume(in.Volume)
	}
	att["machine"] = FlattenMachine(&in.Machine)

	return []interface{}{att}[0]
}
