package shoot

import (
	"testing"

	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	cmp "github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/schema"
)

func getAddonsTestData() (map[string]interface{}, *corev1beta1.Addons) {
	authMode := "auth_mode"
	addonsRaw := map[string]interface{}{
		"kubernetes_dashboard": []interface{}{
			map[string]interface{}{
				"enabled":             true,
				"authentication_mode": authMode,
			},
		},
		"nginx_ingress": []interface{}{
			map[string]interface{}{
				"enabled": true,
			},
		},
	}
	addons := &corev1beta1.Addons{
		KubernetesDashboard: &corev1beta1.KubernetesDashboard{
			Addon: corev1beta1.Addon{
				Enabled: true,
			},
			AuthenticationMode: &authMode,
		},
		NginxIngress: &corev1beta1.NginxIngress{
			Addon: corev1beta1.Addon{
				Enabled: true,
			},
		},
	}

	return addonsRaw, addons
}

func TestExpandAddons(t *testing.T) {
	addonsRaw, addons := getAddonsTestData()
	data := schema.TestResourceDataRaw(t, AddonsResource().Schema, addonsRaw)
	out := ExpandAddons([]interface{}{data.Get("")})
	if diff := cmp.Diff(addons, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}

func TestFlattenAddons(t *testing.T) {
	addonsRaw, addons := getAddonsTestData()
	out := FlattenAddons(addons)[0].(map[string]interface{})
	data := schema.TestResourceDataRaw(t, AddonsResource().Schema, addonsRaw).Get("")
	if diff := cmp.Diff(data, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}
