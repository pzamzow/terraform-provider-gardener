package shoot

import (
	"testing"

	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	cmp "github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/schema"
)

func getKubeletTestData() (map[string]interface{}, *corev1beta1.KubeletConfig) {
	quota := true
	limit := int64(1)
	policy := "policy"
	kubeletRaw := map[string]interface{}{
		"pod_pids_limit":     1,
		"cpu_cfs_quota":      quota,
		"cpu_manager_policy": policy,
	}
	kubelet := &corev1beta1.KubeletConfig{
		PodPIDsLimit:     &limit,
		CPUCFSQuota:      &quota,
		CPUManagerPolicy: &policy,
	}

	return kubeletRaw, kubelet
}

func TestExpandKubelet(t *testing.T) {
	kubeletRaw, kubelet := getKubeletTestData()
	data := schema.TestResourceDataRaw(t, KubeletResource().Schema, kubeletRaw)
	out := ExpandKubelet([]interface{}{data.Get("")})
	if diff := cmp.Diff(kubelet, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}

func TestFlattenKubelet(t *testing.T) {
	kubeletRaw, kubelet := getKubeletTestData()
	out := FlattenKubelet(kubelet)[0].(map[string]interface{})
	data := schema.TestResourceDataRaw(t, KubeletResource().Schema, kubeletRaw).Get("")
	if diff := cmp.Diff(data, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}
