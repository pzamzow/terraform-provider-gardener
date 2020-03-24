package shoot

import (
	"testing"

	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	cmp "github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/schema"
)

func getMonitoringTestData() (map[string]interface{}, *corev1beta1.Monitoring) {
	monitoringRaw := map[string]interface{}{
		"alerting": []interface{}{
			map[string]interface{}{
				"emailreceivers": []interface{}{"receiver1", "receiver2"},
			},
		},
	}
	monitoring := &corev1beta1.Monitoring{
		Alerting: &corev1beta1.Alerting{
			EmailReceivers: []string{"receiver1", "receiver2"},
		},
	}

	return monitoringRaw, monitoring
}

func TestExpandMonitoring(t *testing.T) {
	monitoringRaw, monitoring := getMonitoringTestData()
	data := schema.TestResourceDataRaw(t, MonitoringResource().Schema, monitoringRaw)
	out := ExpandMonitoring([]interface{}{data.Get("")})
	if diff := cmp.Diff(monitoring, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}

func TestFlattenMonitoring(t *testing.T) {
	monitoringRaw, monitoring := getMonitoringTestData()
	out := FlattenMonitoring(monitoring)[0].(map[string]interface{})
	data := schema.TestResourceDataRaw(t, MonitoringResource().Schema, monitoringRaw).Get("")
	if diff := cmp.Diff(data, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}
