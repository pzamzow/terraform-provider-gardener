package gardener

import (
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
)

func ShootResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"addons": {
				Type:        schema.TypeList,
				Description: "Addons contains information about enabled/disabled addons and their configuration.",
				Optional:    true,
				MaxItems:    1,
				Elem:        AddonsResource(),
			},
			"cloud_profile_name": {
				Type:        schema.TypeString,
				Description: "Profile is a name of a CloudProfile object.",
				Required:    true,
			},
			"dns": {
				Type:        schema.TypeList,
				Description: "DNS contains information about the DNS settings of the Shoot.",
				Optional:    true,
				MaxItems:    1,
				Elem:        DNSResource(),
			},
			"hibernation": {
				Type:        schema.TypeList,
				Description: "Hibernation contains information whether the Shoot is suspended or not.",
				Optional:    true,
				MaxItems:    1,
				Elem:        HibernationResource(),
			},
			"kubernetes": {
				Type:        schema.TypeList,
				Description: "Kubernetes contains the version and configuration settings of the control plane components.",
				Required:    true,
				MaxItems:    1,
				Elem:        KubernetesResource(),
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
				DiffSuppressFunc: suppressMissingOptionalConfigurationBlock,
			},
			"monitoring": {
				Type:             schema.TypeList,
				Description:      "Alert configuration to send notification to email lists.",
				Optional:         true,
				MaxItems:         1,
				Elem:             MonitoringResource(),
				DiffSuppressFunc: suppressMissingOptionalConfigurationBlock,
			},
			"provider": {
				Type:             schema.TypeList,
				Description:      "Provider contains all provider-specific and provider-relevant information.",
				Required:         true,
				Elem:             ProviderResource(),
				DiffSuppressFunc: suppressMissingOptionalConfigurationBlock,
			},
			"purpose": {
				Type:        schema.TypeString,
				Description: "Purpose is the purpose class for this cluster.",
				Optional:    true,
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
				Type:        schema.TypeString,
				Description: "Seed name is the name of the seed cluster that runs the control plane of the Shoot.",
				Optional:    true,
			},
		},
	}
}

func ExpandShoot(shoot []interface{}) corev1beta1.ShootSpec {
	obj := corev1beta1.ShootSpec{}

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
		obj.Provider = ExpandProvider(v)
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
		obj.Kubernetes = ExpandKubernetes(v)
	}
	if v, ok := in["maintenance"].([]interface{}); ok && len(v) > 0 {
		obj.Maintenance = ExpandMaintenance(v)
	}

	if v, ok := in["monitoring"].([]interface{}); ok && len(v) > 0 {
		obj.Monitoring = ExpandMonitoring(v)
	}

	if v, ok := in["networking"].([]interface{}); ok && len(v) > 0 {
		obj.Networking = ExpandNetworking(v)
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

func FlattenShoot(in corev1beta1.ShootSpec, d *schema.ResourceData, specPrefix ...string) ([]interface{}, error) {
	att := make(map[string]interface{})
	prefix := ""
	if len(specPrefix) > 0 {
		prefix = specPrefix[0]
	}

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
		att["purpose"] = *in.Purpose
	}
	if in.Addons != nil {
		configAddons := d.Get(prefix + "spec.0.addons").([]interface{})
		flattenedAddons := FlattenAddons(in.Addons)
		att["addons"] = RemoveInternalKeysArraySpec([]interface{}{flattenedAddons}, configAddons)
	}
	configProvider := d.Get(prefix + "spec.0.provider").([]interface{})
	flattenedProvider := FlattenProvider(in.Provider)
	att["provider"] = RemoveInternalKeysArraySpec(flattenedProvider, configProvider)
	if in.DNS != nil {
		configDNS := d.Get(prefix + "spec.0.dns").([]interface{})
		flattenedDNS := FlattenDNS(in.DNS)
		att["dns"] = RemoveInternalKeysArraySpec(flattenedDNS, configDNS)
	}
	if in.Hibernation != nil {
		configHibernation := d.Get(prefix + "spec.0.hibernation").([]interface{})
		flattenedHibernation := FlattenHibernation(in.Hibernation)
		att["hibernation"] = RemoveInternalKeysArraySpec(flattenedHibernation, configHibernation)
	}
	configKubernetes := d.Get(prefix + "spec.0.kubernetes").([]interface{})
	flattenedKubernetes := FlattenKubernetes(in.Kubernetes)
	att["kubernetes"] = RemoveInternalKeysArraySpec(flattenedKubernetes, configKubernetes)
	if in.Maintenance != nil {
		configMaintenance := d.Get(prefix + "spec.0.maintenance").([]interface{})
		flattenedMaintenance := FlattenMaintenance(in.Maintenance)
		att["maintenance"] = RemoveInternalKeysArraySpec(flattenedMaintenance, configMaintenance)
	}
	configNetworking := d.Get(prefix + "spec.0.networking").([]interface{})
	flattenedNetworking := FlattenNetworking(in.Networking)
	att["networking"] = RemoveInternalKeysArraySpec(flattenedNetworking, configNetworking)
	if in.Monitoring != nil {
		configMonitoring := d.Get(prefix + "spec.0.monitoring").([]interface{})
		flattenedMonitoring := FlattenMonitoring(in.Monitoring)
		att["monitoring"] = RemoveInternalKeysArraySpec(flattenedMonitoring, configMonitoring)
	}

	return []interface{}{att}, nil
}
