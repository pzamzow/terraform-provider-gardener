package gardener

import (
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
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
				Type:        schema.TypeList,
				Description: "Kubernetes contains configuration for Kubernetes components related to this worker pool.",
				Optional:    true,
				Elem:        WorkerKubernetesResource(),
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
				Type:        schema.TypeList,
				Description: "MachineType is the machine type of the worker group.",
				Required:    true,
				MaxItems:    1,
				Elem: MachineResource(),
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
				Type:        schema.TypeList,
				Description: "Taints is a list of taints for all the `Node` objects in this worker pool.",
				Optional:    true,
				Elem: TaintResource(),
			},
			"volume": {
				Type:        schema.TypeList,
				Description: "Volume contains information about the volume type and size.",
				Required:    true,
				MaxItems:    1,
				Elem: VolumeResource(),
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

func ExpandWorker(w interface{}) corev1beta1.Worker {
	obj := corev1beta1.Worker{}
	if w == nil {
		return obj
	}
	in := w.(map[string]interface{})

	if v, ok := in["annotations"].(map[string]interface{}); ok {
		obj.Annotations = expandStringMap(v)
	}

	if v, ok := in["cabundle"].(string); ok && len(v) > 0 {
		obj.CABundle = &v
	}

	if v, ok := in["kubernetes"].([]interface{}); ok && len(v) > 0 {
		obj.Kubernetes = ExpandWorkerKubernetes(v)
	}

	if v, ok := in["labels"].(map[string]interface{}); ok {
		obj.Labels = expandStringMap(v)
	}

	if v, ok := in["name"].(string); ok && len(v) > 0 {
		obj.Name = v
	}

	if v, ok := in["machine"].([]interface{}); ok && len(v) > 0 {
		obj.Machine = ExpandMachine(v)
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
		obj.Zones = expandSet(v)
	}

	return obj
}


func FlattenWorker(in []corev1beta1.Worker) []interface{} {
	att := map[string]interface{}{}

	workers := make([]interface{}, len(in))
	for i, v := range in {
		m := map[string]interface{}{}

		if len(v.Name) > 0 {
			m["name"] = v.Name
		}
		if len(v.Zones) > 0 {
			m["zones"] = v.Zones
		}
		if len(v.Taints) > 0 {
			m["taints"] = FlattenTaint(v.Taints)
		}
		if v.MaxSurge != nil {
			m["max_surge"] = v.MaxSurge.IntValue()
		}
		if v.MaxUnavailable != nil {
			m["max_unavailable"] = v.MaxUnavailable.IntValue()
		}
		if v.CABundle != nil {
			m["cabundle"] = *v.CABundle
		}

		if v.Minimum != 0 {
			m["minimum"] = v.Minimum
		}

		if v.Maximum != 0 {
			m["maximum"] = v.Maximum
		}

		if v.Kubernetes != nil {
			m["kubernetes"] = FlattenWorkerKubernetes(v.Kubernetes)
		}

		if len(v.Annotations) > 0 {
			m["annotations"] = v.Annotations
		}
		if len(v.Labels) > 0 {
			m["labels"] = v.Labels
		}
		if v.Volume != nil {
			m["volume"] = FlattenVolume(v.Volume)
		}
		m["machine"] = FlattenMachine(v.Machine)

		workers[i] = m
	}
	att["worker"] = workers

	return []interface{}{att}
}