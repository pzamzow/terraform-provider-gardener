package gardener

import (
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
)

func KubernetesResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"allow_privileged_containers": {
				Type:        schema.TypeBool,
				Description: "AllowPrivilegedContainers indicates whether privileged containers are allowed in the Shoot.",
				Optional:    true,
			},
			"kube_api_server": {
				Type:        schema.TypeList,
				Description: "KubeAPIServer contains configuration settings for the kube-apiserver.",
				Optional:    true,
				MaxItems:    1,
				Elem: KubeAPIServerResource(),
			},
			"kube_controller_manager": {
				Type:        schema.TypeList,
				Description: "KubeControllerManager contains configuration settings for the kube-controller-manager.",
				Optional:    true,
				MaxItems:    1,
				Elem: KubeControllerManagerResource(),
			},
			"kube_proxy": {
				Type:        schema.TypeList,
				Description: "KubeProxy contains configuration settings for the kube-proxy.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:        schema.TypeString,
							Description: "Mode specifies which proxy mode to use. defaults to IPTables.",
							Optional:    true,

							Default: "IPTables",
						},
					},
				},
			},
			"kubelet": {
				Type:        schema.TypeList,
				Description: "Kubelet contains configuration settings for the kubelet.",
				Optional:    true,
				MaxItems:    1,
				Elem: KubeletResource(),
			},
			"version": {
				Type:        schema.TypeString,
				Description: "Version is the semantic Kubernetes version to use for the Shoot cluster.",
				Required:    true,
			},
			"cluster_autoscaler": {
				Type:        schema.TypeList,
				Description: "ClusterAutoscaler contains the configration flags for the Kubernetes cluster autoscaler.",
				Optional:    true,
				MaxItems:    1,
				Elem: ClusterAutoscalerResource(),
			},
		},
	}
}

func ExpandKubernetes(kubernetes []interface{}) corev1beta1.Kubernetes {
	obj := corev1beta1.Kubernetes{}

	if len(kubernetes) == 0 || kubernetes[0] == nil {
		return obj
	}
	in := kubernetes[0].(map[string]interface{})

	if v, ok := in["allow_privileged_containers"].(bool); ok {
		obj.AllowPrivilegedContainers = &v
	}
	if v, ok := in["kube_api_server"].([]interface{}); ok && len(v) > 0 {
		obj.KubeAPIServer = ExpandKubeAPIServer(v)
	}
	if v, ok := in["kube_controller_manager"].([]interface{}); ok && len(v) > 0 {
		obj.KubeControllerManager = ExpandKubeControllerManager(v)
	}
	if v, ok := in["kube_scheduler"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.KubeScheduler = &corev1beta1.KubeSchedulerConfig{}

		if v, ok := v["feature_gates"].(map[string]interface{}); ok {
			obj.KubeScheduler.FeatureGates = expandBoolMap(v)
		}
	}
	if v, ok := in["kube_proxy"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.KubeProxy = &corev1beta1.KubeProxyConfig{}

		if v, ok := v["feature_gates"].(map[string]interface{}); ok {
			obj.KubeProxy.FeatureGates = expandBoolMap(v)
		}
		if v, ok := v["mode"].(string); ok && len(v) > 0 {
			mode := corev1beta1.ProxyModeIPTables
			obj.KubeProxy.Mode = &mode
		}
	}
	if v, ok := in["kubelet"].([]interface{}); ok && len(v) > 0 {
		obj.Kubelet = ExpandKubelet(v)
	}
	if v, ok := in["version"].(string); ok {
		obj.Version = v
	}
	if v, ok := in["cluster_autoscaler"].([]interface{}); ok && len(v) > 0 {
		obj.ClusterAutoscaler = ExpandClusterAutoscaler(v)
	}

	return obj
}

func FlattenKubernetes(in corev1beta1.Kubernetes) []interface{} {
	att := make(map[string]interface{})

	if in.AllowPrivilegedContainers != nil {
		att["allow_privileged_containers"] = in.AllowPrivilegedContainers
	}
	if in.KubeAPIServer != nil {
		att["kube_api_server"] = FlattenKubeAPIServer(in.KubeAPIServer)
	}
	if in.KubeControllerManager != nil {
		att["kube_controller_manager"] = FlattenKubeControllerManager(in.KubeControllerManager)
	}
	if in.KubeProxy != nil {
		proxy := make(map[string]interface{})
		if in.KubeProxy.FeatureGates != nil {
			proxy["feature_gates"] = in.KubeProxy.FeatureGates
		}
		if in.KubeProxy.Mode != nil {
			proxy["mode"] = in.KubeProxy.Mode
		}
		att["kube_proxy"] = []interface{}{proxy}
	}
	if in.Kubelet != nil {
		att["kubelet"] = FlattenKubelet(in.Kubelet)
	}
	if len(in.Version) > 0 {
		att["version"] = in.Version
	}
	if in.ClusterAutoscaler != nil {
		att["cluster_autoscaler"] = FlattenClusterAutoscaler(in.ClusterAutoscaler)
	}

	return []interface{}{att}
}
