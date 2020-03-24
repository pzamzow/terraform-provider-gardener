package shoot

import (
	"testing"

	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	cmp "github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/schema"
)

func getWorkerKubernetesTestData() (map[string]interface{}, *corev1beta1.WorkerKubernetes) {
	kubeletRaw, kubelet := getKubeletTestData()
	workerKubernetesRaw := map[string]interface{}{
		"kubelet": []interface{}{kubeletRaw},
	}
	workerKubernetes := &corev1beta1.WorkerKubernetes{
		Kubelet: kubelet,
	}

	return workerKubernetesRaw, workerKubernetes
}

func TestExpandWorkerKubernetes(t *testing.T) {
	workerKubernetesRaw, workerKubernetes := getWorkerKubernetesTestData()
	data := schema.TestResourceDataRaw(t, WorkerKubernetesResource().Schema, workerKubernetesRaw)
	out := ExpandWorkerKubernetes([]interface{}{data.Get("")})
	if diff := cmp.Diff(workerKubernetes, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}

func TestFlattenWorkerKubernetes(t *testing.T) {
	workerKubernetesRaw, workerKubernetes := getWorkerKubernetesTestData()
	out := FlattenWorkerKubernetes(workerKubernetes)[0].(map[string]interface{})
	data := schema.TestResourceDataRaw(t, WorkerKubernetesResource().Schema, workerKubernetesRaw).Get("")
	if diff := cmp.Diff(data, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}
