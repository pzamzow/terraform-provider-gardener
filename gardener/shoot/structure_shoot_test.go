package shoot

import (
	"testing"

	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	cmp "github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/schema"
)

func getShootTestData() (map[string]interface{}, *corev1beta1.ShootSpec) {
	addonsRaw, addons := getAddonsTestData()
	dnsRaw, dns := getDNSTestData()
	hibernationRaw, hibernation := getHibernationTestData()
	kubernetesRaw, kubernetes := getKubernetesTestData()
	networkingRaw, networking := getNetworkingTestData()
	maintenanceRaw, maintenance := getMaintenanceTestData()
	monitoringRaw, monitoring := getMonitoringTestData()
	providerRaw, provider := getProviderTestData("azure")
	purpose := corev1beta1.ShootPurposeProduction
	seed := "seed"
	shootRaw := map[string]interface{}{
		"cloud_profile_name":  "azure",
		"addons":              []interface{}{addonsRaw},
		"dns":                 []interface{}{dnsRaw},
		"hibernation":         []interface{}{hibernationRaw},
		"kubernetes":          []interface{}{kubernetesRaw},
		"networking":          []interface{}{networkingRaw},
		"maintenance":         []interface{}{maintenanceRaw},
		"monitoring":          []interface{}{monitoringRaw},
		"provider":            []interface{}{providerRaw},
		"purpose":             string(purpose),
		"region":              "region",
		"secret_binding_name": "secret",
		"seed_name":           seed,
	}
	shoot := &corev1beta1.ShootSpec{
		CloudProfileName:  "azure",
		Addons:            addons,
		DNS:               dns,
		Hibernation:       hibernation,
		Kubernetes:        *kubernetes,
		Networking:        *networking,
		Maintenance:       maintenance,
		Monitoring:        monitoring,
		Provider:          *provider,
		Purpose:           &purpose,
		Region:            "region",
		SecretBindingName: "secret",
		SeedName:          &seed,
	}

	return shootRaw, shoot
}

func TestExpandShoot(t *testing.T) {
	shootRaw, shoot := getShootTestData()
	data := schema.TestResourceDataRaw(t, ShootResource().Schema, shootRaw)
	out := ExpandShoot([]interface{}{data.Get("")})
	if diff := cmp.Diff(shoot, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}

func TestFlattenShoot(t *testing.T) {
	shootRaw, shoot := getShootTestData()
	out := FlattenShoot(shoot)
	data := schema.TestResourceDataRaw(t, ShootResource().Schema, shootRaw)
	if diff := cmp.Diff(data.Get(""), out[0].(map[string]interface{})); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}
