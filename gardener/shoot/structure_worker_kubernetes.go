package shoot

import (
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
	lib "github.com/kyma-incubator/terraform-provider-gardener/gardener/lib"
)

func WorkerKubernetesResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"kubelet": {
				Type:             schema.TypeList,
				Description:      "Kubelet contains configuration settings for the kubelet.",
				Optional:         true,
				MaxItems:         1,
				Elem:             KubeletResource(),
				DiffSuppressFunc: lib.SuppressMissingOptionalConfigurationBlock,
			},
		},
	}
}

func ExpandWorkerKubernetes(wk []interface{}) *corev1beta1.WorkerKubernetes {
	obj := &corev1beta1.WorkerKubernetes{}
	if len(wk) == 0 && wk[0] == nil {
		return obj
	}
	in := wk[0].(map[string]interface{})

	if v, ok := in["kubelet"].([]interface{}); ok {
		obj.Kubelet = ExpandKubelet(v)
	}

	return obj
}

func FlattenWorkerKubernetes(in *corev1beta1.WorkerKubernetes) []interface{} {
	att := make(map[string]interface{})

	if in.Kubelet != nil {
		att["kubelet"] = FlattenKubelet(in.Kubelet)
	}

	return []interface{}{att}
}
