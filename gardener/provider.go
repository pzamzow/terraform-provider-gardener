package gardener

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	lib "github.com/kyma-incubator/terraform-provider-gardener/gardener/lib"
	shoot "github.com/kyma-incubator/terraform-provider-gardener/gardener/shoot"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"kube_file": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"gardener_shoot": shoot.ResourceShoot(),
		},
		ConfigureFunc: providerConfigure,
	}
}
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := &lib.Config{
		KubeFile: d.Get("kube_file").(string),
	}
	return lib.NewClient(config)
}
