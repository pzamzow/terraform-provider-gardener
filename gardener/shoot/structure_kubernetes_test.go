package shoot

import (
	"testing"

	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	cmp "github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/schema"
)

func getKubernetesTestData() (map[string]interface{}, *corev1beta1.Kubernetes) {
	kubeAPIServerRaw, kubeAPIServer := getKubeAPIServerTestData()
	kubeControllerManagerRaw, kubeControllerManager := getKubeControllerManagerTestData()
	kubeletRaw, kubelet := getKubeletTestData()
	clusterAutoscalerRaw, clusterAutoscaler := getClusterAutoscalerTestData()
	privileged := true
	mode := corev1beta1.ProxyModeIPVS
	kubernetesRaw := map[string]interface{}{
		"allow_privileged_containers": privileged,
		"kube_api_server":             []interface{}{kubeAPIServerRaw},
		"kube_controller_manager":     []interface{}{kubeControllerManagerRaw},
		"kube_proxy": []interface{}{
			map[string]interface{}{
				"mode": "IPVS",
			},
		},
		"kubelet":            []interface{}{kubeletRaw},
		"version":            "1.0.0",
		"cluster_autoscaler": []interface{}{clusterAutoscalerRaw},
	}
	kubernetes := &corev1beta1.Kubernetes{
		AllowPrivilegedContainers: &privileged,
		KubeAPIServer:             kubeAPIServer,
		KubeControllerManager:     kubeControllerManager,
		KubeProxy: &corev1beta1.KubeProxyConfig{
			Mode: &mode,
		},
		Kubelet:           kubelet,
		Version:           "1.0.0",
		ClusterAutoscaler: clusterAutoscaler,
	}

	return kubernetesRaw, kubernetes
}

func TestExpandKubernetes(t *testing.T) {
	kubernetesRaw, kubernetes := getKubernetesTestData()
	data := schema.TestResourceDataRaw(t, KubernetesResource().Schema, kubernetesRaw)
	out := ExpandKubernetes([]interface{}{data.Get("")})
	if diff := cmp.Diff(kubernetes, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}

func TestFlattenKubernetes(t *testing.T) {
	kubernetesRaw, kubernetes := getKubernetesTestData()
	out := FlattenKubernetes(kubernetes)[0].(map[string]interface{})
	data := schema.TestResourceDataRaw(t, KubernetesResource().Schema, kubernetesRaw).Get("")
	if diff := cmp.Diff(data, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}
