package shoot

import (
	"testing"
	"time"

	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	cmp "github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/schema"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getClusterAutoscalerTestData() (map[string]interface{}, *corev1beta1.ClusterAutoscaler) {
	threshold := 1.0
	clusterAutoscalerRaw := map[string]interface{}{
		"scale_down_utilization_threshold": threshold,
		"scale_down_unneeded_time":         "10s",
		"scale_down_delay_after_add":       "20s",
		"scale_down_delay_after_failure":   "30s",
		"scale_down_delay_after_delete":    "40s",
		"scan_interval":                    "50s",
	}
	clusterAutoscaler := &corev1beta1.ClusterAutoscaler{
		ScaleDownUtilizationThreshold: &threshold,
		ScaleDownUnneededTime:         &v1.Duration{Duration: 10 * time.Second},
		ScaleDownDelayAfterAdd:        &v1.Duration{Duration: 20 * time.Second},
		ScaleDownDelayAfterFailure:    &v1.Duration{Duration: 30 * time.Second},
		ScaleDownDelayAfterDelete:     &v1.Duration{Duration: 40 * time.Second},
		ScanInterval:                  &v1.Duration{Duration: 50 * time.Second},
	}

	return clusterAutoscalerRaw, clusterAutoscaler
}

func TestExpandClusterAutoscaler(t *testing.T) {
	clusterAutoscalerRaw, clusterAutoscaler := getClusterAutoscalerTestData()
	data := schema.TestResourceDataRaw(t, ClusterAutoscalerResource().Schema, clusterAutoscalerRaw)
	out := ExpandClusterAutoscaler([]interface{}{data.Get("")})
	if diff := cmp.Diff(clusterAutoscaler, out); diff != "" {
		t.Fatalf("Error clusterAutoScaler output and expected: \n%s", diff)
	}
}

func TestFlattenClusterAutoscaler(t *testing.T) {
	clusterAutoscalerRaw, clusterAutoscaler := getClusterAutoscalerTestData()
	out := FlattenClusterAutoscaler(clusterAutoscaler)[0].(map[string]interface{})
	data := schema.TestResourceDataRaw(t, ClusterAutoscalerResource().Schema, clusterAutoscalerRaw).Get("")
	if diff := cmp.Diff(data, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}
