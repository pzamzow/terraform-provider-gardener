package shoot

import (
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
	lib "github.com/kyma-incubator/terraform-provider-gardener/gardener/lib"
)

func DNSResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:             schema.TypeString,
				Description:      "Domain is the external available domain of the Shoot cluster.",
				Optional:         true,
				DiffSuppressFunc: lib.SuppressEmptyNewValue,
			},
		},
	}
}

func ExpandDNS(dns []interface{}) *corev1beta1.DNS {
	obj := corev1beta1.DNS{}

	if len(dns) == 0 || dns[0] == nil {
		return &obj
	}
	in := dns[0].(map[string]interface{})

	if v, ok := in["domain"].(string); ok && len(v) > 0 {
		obj.Domain = &v
	}

	return &obj
}

func FlattenDNS(in *corev1beta1.DNS) []interface{} {
	att := make(map[string]interface{})

	if in.Domain != nil {
		att["domain"] = *in.Domain
	}

	return []interface{}{att}
}
