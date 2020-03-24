package shoot

import (
	"testing"

	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	cmp "github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func getWorkerTestData() (map[string]interface{}, *corev1beta1.Worker) {
	workerKubernetesRaw, workerKubernetes := getWorkerKubernetesTestData()
	machineRaw, machine := getMachineTestData()
	taintRaw, taint := getTaintTestData()
	volumeRaw, volume := getVolumeTestData()
	caBundle := "caBundle"
	workerRaw := map[string]interface{}{
		"annotations": map[string]interface{}{
			"test-key-annotation": "test-value-annotation",
		},
		"cabundle":   caBundle,
		"kubernetes": []interface{}{workerKubernetesRaw},
		"labels": map[string]interface{}{
			"test-key-label": "test-value-label",
		},
		"name":            "cpu-worker",
		"machine":         []interface{}{machineRaw},
		"minimum":         1,
		"maximum":         2,
		"max_surge":       1,
		"max_unavailable": 0,
		"taints":          taintRaw,
		"volume":          []interface{}{volumeRaw},
		"zones":           []interface{}{"foo", "bar"},
	}
	worker := &corev1beta1.Worker{
		Annotations: map[string]string{
			"test-key-annotation": "test-value-annotation",
		},
		CABundle:   &caBundle,
		Kubernetes: workerKubernetes,
		Labels: map[string]string{
			"test-key-label": "test-value-label",
		},
		Name:    "cpu-worker",
		Machine: *machine,
		Minimum: int32(1),
		Maximum: int32(2),
		MaxSurge: &intstr.IntOrString{
			IntVal: 1,
		},
		MaxUnavailable: &intstr.IntOrString{
			IntVal: 0,
		},
		Taints: taint,
		Volume: volume,
		Zones:  []string{"bar", "foo"},
	}

	return workerRaw, worker
}

func TestExpandWorker(t *testing.T) {
	workerRaw, worker := getWorkerTestData()
	data := schema.TestResourceDataRaw(t, WorkerResource().Schema, workerRaw)
	out := ExpandWorker([]interface{}{data.Get("")}[0])
	if diff := cmp.Diff(worker, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}

func TestFlattenWorker(t *testing.T) {
	workerRaw, worker := getWorkerTestData()
	out := FlattenWorker(worker).(map[string]interface{})
	data := schema.TestResourceDataRaw(t, WorkerResource().Schema, workerRaw).Get("")
	if diff := cmp.Diff(data, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}
