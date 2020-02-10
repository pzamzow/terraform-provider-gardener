package gardener

import (
	"testing"

	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	cmp "github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/schema"
)

func TestExpandWorker(t *testing.T) {
	authMode := "auth_mode"

	addons := map[string]interface{}{
		"addons": []interface{}{
			map[string]interface{}{
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
			},
		},
	}
	expected := &corev1beta1.Addons{
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

	data := schema.TestResourceDataRaw(t, AddonsResource().Schema, addons)
	out := ExpandAddons(data.Get("addons").([]interface{}))
	if diff := cmp.Diff(expected, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}

func TestFlattenWorker(t *testing.T) {
	authMode := "auth_mode"

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
	expected := []interface{}{
		map[string]interface{}{
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
		},
	}

	out := FlattenAddons(addons)
	data := schema.TestResourceDataRaw(t, AddonsResource().Schema, expected)
	if diff := cmp.Diff(data, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}
