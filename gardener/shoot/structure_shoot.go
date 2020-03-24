package shoot

import (
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
	lib "github.com/kyma-incubator/terraform-provider-gardener/gardener/lib"
)

func ShootResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"addons": {
				Type:             schema.TypeList,
				Description:      "Addons contains information about enabled/disabled addons and their configuration.",
				Optional:         true,
				MaxItems:         1,
				Elem:             AddonsResource(),
				DiffSuppressFunc: lib.SuppressMissingOptionalConfigurationBlock,
			},
			"cloud_profile_name": {
				Type:        schema.TypeString,
				Description: "Profile is a name of a CloudProfile object.",
				Required:    true,
			},
			"dns": {
				Type:             schema.TypeList,
				Description:      "DNS contains information about the DNS settings of the Shoot.",
				Optional:         true,
				MaxItems:         1,
				Elem:             DNSResource(),
				DiffSuppressFunc: lib.SuppressMissingOptionalConfigurationBlock,
			},
			"hibernation": {
				Type:             schema.TypeList,
				Description:      "Hibernation contains information whether the Shoot is suspended or not.",
				Optional:         true,
				MaxItems:         1,
				Elem:             HibernationResource(),
				DiffSuppressFunc: lib.SuppressMissingOptionalConfigurationBlock,
			},
			"kubernetes": {
				Type:             schema.TypeList,
				Description:      "Kubernetes contains the version and configuration settings of the control plane components.",
				Required:         true,
				MaxItems:         1,
				Elem:             KubernetesResource(),
				DiffSuppressFunc: lib.SuppressMissingOptionalConfigurationBlock,
			},
			"networking": {
				Type:        schema.TypeList,
				Description: "Networking contains information about cluster networking such as CNI Plugin type, CIDRs, ...etc.",
				Required:    true,
				Elem:        NetworkingResource(),
			},
			"maintenance": {
				Type:             schema.TypeList,
				Description:      "Maintenance contains information about the time window for maintenance operations and which operations should be performed.",
				Optional:         true,
				MaxItems:         1,
				Elem:             MaintenanceResource(),
				DiffSuppressFunc: lib.SuppressMissingOptionalConfigurationBlock,
			},
			"monitoring": {
				Type:             schema.TypeList,
				Description:      "Alert configuration to send notification to email lists.",
				Optional:         true,
				MaxItems:         1,
				Elem:             MonitoringResource(),
				DiffSuppressFunc: lib.SuppressMissingOptionalConfigurationBlock,
			},
			"provider": {
				Type:             schema.TypeList,
				Description:      "Provider contains all provider-specific and provider-relevant information.",
				Required:         true,
				Elem:             ProviderResource(),
				DiffSuppressFunc: lib.SuppressMissingOptionalConfigurationBlock,
			},
			"purpose": {
				Type:             schema.TypeString,
				Description:      "Purpose is the purpose class for this cluster.",
				Optional:         true,
				DiffSuppressFunc: lib.SuppressEmptyNewValue,
			},
			"region": {
				Type:        schema.TypeString,
				Description: "Region is a name of a cloud provider region.",
				Required:    true,
			},
			"secret_binding_name": {
				Type:        schema.TypeString,
				Description: "Secret binding name is the name of the a SecretBinding that has a reference to the provider secret. The credentials inside the provider secret will be used to create the shoot in the respective account",
				Required:    true,
			},
			"seed_name": {
				Type:             schema.TypeString,
				Description:      "Seed name is the name of the seed cluster that runs the control plane of the Shoot.",
				Optional:         true,
				DiffSuppressFunc: lib.SuppressEmptyNewValue,
			},
		},
	}
}

func ExpandShoot(shoot []interface{}) *corev1beta1.ShootSpec {
	obj := &corev1beta1.ShootSpec{}

	if len(shoot) == 0 || shoot[0] == nil {
		return obj
	}
	in := shoot[0].(map[string]interface{})

	if v, ok := in["addons"].([]interface{}); ok && len(v) > 0 {
		obj.Addons = ExpandAddons(v)
	}
	if v, ok := in["cloud_profile_name"].(string); ok && len(v) > 0 {
		obj.CloudProfileName = v
	}
	if v, ok := in["purpose"].(string); ok && len(v) > 0 {
		purpose := corev1beta1.ShootPurpose(v)
		obj.Purpose = &purpose
	}
	if v, ok := in["provider"].([]interface{}); ok && len(v) > 0 {
		obj.Provider = *ExpandProvider(v)
	}

	if v, ok := in["dns"].([]interface{}); ok && len(v) > 0 {
		dns := ExpandDNS(v)
		if dns.Domain != nil {
			obj.DNS = dns
		}
	}
	if v, ok := in["hibernation"].([]interface{}); ok && len(v) > 0 {
		obj.Hibernation = ExpandHibernation(v)
	}
	if v, ok := in["kubernetes"].([]interface{}); ok && len(v) > 0 {
		obj.Kubernetes = *ExpandKubernetes(v)
	}
	if v, ok := in["maintenance"].([]interface{}); ok && len(v) > 0 {
		obj.Maintenance = ExpandMaintenance(v)
	}

	if v, ok := in["monitoring"].([]interface{}); ok && len(v) > 0 {
		obj.Monitoring = ExpandMonitoring(v)
	}

	if v, ok := in["networking"].([]interface{}); ok && len(v) > 0 {
		obj.Networking = *ExpandNetworking(v)
	}

	if v, ok := in["region"].(string); ok && len(v) > 0 {
		obj.Region = v
	}

	if v, ok := in["secret_binding_name"].(string); ok && len(v) > 0 {
		obj.SecretBindingName = v
	}

	if v, ok := in["seed_name"].(string); ok && len(v) > 0 {
		obj.SeedName = &v
	}

	return obj
}

func FlattenShoot(in *corev1beta1.ShootSpec) []interface{} {
	att := make(map[string]interface{})

	if len(in.CloudProfileName) > 0 {
		att["cloud_profile_name"] = in.CloudProfileName
	}
	if len(in.SecretBindingName) > 0 {
		att["secret_binding_name"] = in.SecretBindingName
	}
	if len(in.Region) > 0 {
		att["region"] = in.Region
	}
	if in.Purpose != nil {
		att["purpose"] = string(*in.Purpose)
	}
	if in.SeedName != nil {
		att["seed_name"] = *in.SeedName
	}
	if in.Addons != nil {
		att["addons"] = FlattenAddons(in.Addons)
	}
	att["provider"] = FlattenProvider(&in.Provider)
	if in.DNS != nil {
		att["dns"] = FlattenDNS(in.DNS)
	}
	if in.Hibernation != nil {
		att["hibernation"] = FlattenHibernation(in.Hibernation)
	}
	att["kubernetes"] = FlattenKubernetes(&in.Kubernetes)
	if in.Maintenance != nil {
		att["maintenance"] = FlattenMaintenance(in.Maintenance)
	}
	att["networking"] = FlattenNetworking(&in.Networking)
	if in.Monitoring != nil {
		att["monitoring"] = FlattenMonitoring(in.Monitoring)
	}

	return []interface{}{att}
}
