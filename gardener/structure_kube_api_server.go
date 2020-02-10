package gardener

import (
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
	v1 "k8s.io/api/core/v1"
)

func KubeAPIServerResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enable_basic_authentication": {
				Type:        schema.TypeBool,
				Description: "enable basic authentication flag.",
				Optional:    true,
			},
			"oidc_config": {
				Type:             schema.TypeList,
				Description:      "interface for adding oidc_config in kube api server section",
				MaxItems:         1,
				Optional:         true,
				DiffSuppressFunc: suppressMissingOptionalConfigurationBlock,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ca_bundle": {
							Type:        schema.TypeString,
							Description: "ca_bundle for oidc config in kube api server section",
							Optional:    true,
						},
						"client_id": {
							Type:        schema.TypeString,
							Description: "client_id for oidc config in kube api server section",
							Optional:    true,
						},
						"groups_claim": {
							Type:        schema.TypeString,
							Description: "groups_claim for oidc config in kube api server section",
							Optional:    true,
						},
						"groups_prefix": {
							Type:        schema.TypeString,
							Description: "groups_prefix for oidc config in kube api server section",
							Optional:    true,
						},
						"issuer_url": {
							Type:        schema.TypeString,
							Description: "issuer_url for oidc config in kube api server section",
							Optional:    true,
						},
						"required_claims": {
							Type:        schema.TypeString,
							Description: "required_claims for oidc config in kube api server section",
							Optional:    true,
						},
						"signing_algs": {
							Type:        schema.TypeString,
							Description: "signing_algs for oidc config in kube api server section",
							Optional:    true,
						},
						"username_claim": {
							Type:        schema.TypeString,
							Description: "username_claim for oidc config in kube api server section",
							Optional:    true,
						},
						"username_prefix": {
							Type:        schema.TypeString,
							Description: "username_prefix for oidc config in kube api server section",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}
func ExpandKubeAPIServer(server []interface{}) *corev1beta1.KubeAPIServerConfig {
	obj := &corev1beta1.KubeAPIServerConfig{}

	if len(server) == 0 || server[0] == nil {
		return obj
	}
	in := server[0].(map[string]interface{})

	if v, ok := in["feature_gates"].(map[string]interface{}); ok {
		obj.FeatureGates = expandBoolMap(v)
	}
	if v, ok := in["runtime_config"].(map[string]interface{}); ok {
		obj.RuntimeConfig = expandBoolMap(v)
	}
	if v, ok := in["oidc_config"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.OIDCConfig = &corev1beta1.OIDCConfig{}

		if v, ok := v["ca_bundle"].(string); ok && len(v) > 0 {
			obj.OIDCConfig.CABundle = &v
		}
		if v, ok := v["client_id"].(string); ok && len(v) > 0 {
			obj.OIDCConfig.ClientID = &v
		}
		if v, ok := v["groups_claim"].(string); ok && len(v) > 0 {
			obj.OIDCConfig.GroupsClaim = &v
		}
		if v, ok := v["groups_prefix"].(string); ok && len(v) > 0 {
			obj.OIDCConfig.GroupsPrefix = &v
		}
		if v, ok := v["issuer_url"].(string); ok && len(v) > 0 {
			obj.OIDCConfig.IssuerURL = &v
		}
		if v, ok := v["required_claims"].(map[string]interface{}); ok {
			obj.OIDCConfig.RequiredClaims = expandStringMap(v)
		}
		if v, ok := v["signing_algs"].(*schema.Set); ok {
			obj.OIDCConfig.SigningAlgs = expandSet(v)
		}
		if v, ok := v["username_claim"].(string); ok && len(v) > 0 {
			obj.OIDCConfig.UsernameClaim = &v
		}
		if v, ok := v["username_prefix"].(string); ok && len(v) > 0 {
			obj.OIDCConfig.UsernamePrefix = &v
		}
	}
	if v, ok := in["audit_config"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.AuditConfig = &corev1beta1.AuditConfig{}

		if v, ok := v["audit_policy"].([]interface{}); ok && len(v) > 0 {
			v := v[0].(map[string]interface{})
			obj.AuditConfig.AuditPolicy = &corev1beta1.AuditPolicy{}

			if v, ok := v["config_map_ref"].([]interface{}); ok {
				obj.AuditConfig.AuditPolicy.ConfigMapRef = expandObjectReference(v)
			}
		}
	}

	return obj
}

func expandObjectReference(l []interface{}) *v1.ObjectReference {
	if len(l) == 0 || l[0] == nil {
		return &v1.ObjectReference{}
	}
	in := l[0].(map[string]interface{})
	obj := &v1.ObjectReference{}
	if v, ok := in["name"].(string); ok {
		obj.Name = v
	}
	return obj
}

func FlattenKubeAPIServer(in *corev1beta1.KubeAPIServerConfig) []interface{} {
	att := make(map[string]interface{})

	if in.EnableBasicAuthentication != nil {
		att["enable_basic_authentication"] = *in.EnableBasicAuthentication
	}

	return []interface{}{att}
}
