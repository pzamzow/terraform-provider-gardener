package shoot

import (
	"testing"

	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	cmp "github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/schema"
)

func getKubeControllerManagerTestData() (map[string]interface{}, *corev1beta1.KubeControllerManagerConfig) {
	size := int32(1)
	kubeControllerManagerRaw := map[string]interface{}{
		"node_cidr_mask_size": 1,
	}
	kubeControllerManager := &corev1beta1.KubeControllerManagerConfig{
		NodeCIDRMaskSize: &size,
	}

	return kubeControllerManagerRaw, kubeControllerManager
}

func TestExpandKubeControllerManager(t *testing.T) {
	kubeControllerManagerRaw, kubeControllerManager := getKubeControllerManagerTestData()
	data := schema.TestResourceDataRaw(t, KubeControllerManagerResource().Schema, kubeControllerManagerRaw)
	out := ExpandKubeControllerManager([]interface{}{data.Get("")})
	if diff := cmp.Diff(kubeControllerManager, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}

func TestFlattenKubeControllerManager(t *testing.T) {
	kubeControllerManagerRaw, kubeControllerManager := getKubeControllerManagerTestData()
	out := FlattenKubeControllerManager(kubeControllerManager)[0].(map[string]interface{})
	data := schema.TestResourceDataRaw(t, KubeControllerManagerResource().Schema, kubeControllerManagerRaw).Get("")
	if diff := cmp.Diff(data, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}
