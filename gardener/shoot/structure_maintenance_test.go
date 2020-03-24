package shoot

import (
	"testing"

	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	cmp "github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/schema"
)

func getMaintenanceTestData() (map[string]interface{}, *corev1beta1.Maintenance) {
	maintenanceRaw := map[string]interface{}{
		"auto_update": []interface{}{
			map[string]interface{}{
				"kubernetes_version":    true,
				"machine_image_version": true,
			},
		},
		"time_window": []interface{}{
			map[string]interface{}{
				"begin": "030000+0000",
				"end":   "040000+0000",
			},
		},
	}
	maintenance := &corev1beta1.Maintenance{
		AutoUpdate: &corev1beta1.MaintenanceAutoUpdate{
			KubernetesVersion:   true,
			MachineImageVersion: true,
		},
		TimeWindow: &corev1beta1.MaintenanceTimeWindow{
			Begin: "030000+0000",
			End:   "040000+0000",
		},
	}

	return maintenanceRaw, maintenance
}

func TestExpandMaintenance(t *testing.T) {
	maintenanceRaw, maintenance := getMaintenanceTestData()
	data := schema.TestResourceDataRaw(t, MaintenanceResource().Schema, maintenanceRaw)
	out := ExpandMaintenance([]interface{}{data.Get("")})
	if diff := cmp.Diff(maintenance, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}

func TestFlattenMaintenance(t *testing.T) {
	maintenanceRaw, maintenance := getMaintenanceTestData()
	out := FlattenMaintenance(maintenance)[0].(map[string]interface{})
	data := schema.TestResourceDataRaw(t, MaintenanceResource().Schema, maintenanceRaw).Get("")
	if diff := cmp.Diff(data, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}
