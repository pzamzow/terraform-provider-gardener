package shoot

import (
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
	lib "github.com/kyma-incubator/terraform-provider-gardener/gardener/lib"
)

func KubeAPIServerResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enable_basic_authentication": {
				Type:        schema.TypeBool,
				Description: "enable basic authentication flag.",
				Optional:    true,
				Default:     false,
			},
			"oidc_config": {
				Type:             schema.TypeList,
				Description:      "interface for adding oidc_config in kube api server section",
				MaxItems:         1,
				Optional:         true,
				DiffSuppressFunc: lib.SuppressMissingOptionalConfigurationBlock,
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
							Type:        schema.TypeMap,
							Description: "required_claims for oidc config in kube api server section",
							Optional:    true,
						},
						"signing_algs": {
							Type:        schema.TypeSet,
							Description: "signing_algs for oidc config in kube api server section",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         schema.HashString,
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

	if v, ok := in["enable_basic_authentication"].(bool); ok {
		obj.EnableBasicAuthentication = &v
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
			obj.OIDCConfig.RequiredClaims = lib.ExpandStringMap(v)
		}
		if v, ok := v["signing_algs"].(*schema.Set); ok {
			obj.OIDCConfig.SigningAlgs = lib.ExpandSet(v)
		}
		if v, ok := v["username_claim"].(string); ok && len(v) > 0 {
			obj.OIDCConfig.UsernameClaim = &v
		}
		if v, ok := v["username_prefix"].(string); ok && len(v) > 0 {
			obj.OIDCConfig.UsernamePrefix = &v
		}
	}

	return obj
}

func FlattenKubeAPIServer(in *corev1beta1.KubeAPIServerConfig) []interface{} {
	att := make(map[string]interface{})

	if in.EnableBasicAuthentication != nil {
		att["enable_basic_authentication"] = *in.EnableBasicAuthentication
	}

	if in.OIDCConfig != nil {
		config := make(map[string]interface{})

		if in.OIDCConfig.CABundle != nil {
			config["ca_bundle"] = *in.OIDCConfig.CABundle
		}
		if in.OIDCConfig.ClientID != nil {
			config["client_id"] = *in.OIDCConfig.ClientID
		}
		if in.OIDCConfig.GroupsClaim != nil {
			config["groups_claim"] = *in.OIDCConfig.GroupsClaim
		}
		if in.OIDCConfig.GroupsPrefix != nil {
			config["groups_prefix"] = *in.OIDCConfig.GroupsPrefix
		}
		if in.OIDCConfig.IssuerURL != nil {
			config["issuer_url"] = *in.OIDCConfig.IssuerURL
		}
		if in.OIDCConfig.RequiredClaims != nil {
			config["required_claims"] = lib.FlattenStringMap(in.OIDCConfig.RequiredClaims)
		}
		if in.OIDCConfig.SigningAlgs != nil {
			config["signing_algs"] = lib.NewStringSet(schema.HashString, in.OIDCConfig.SigningAlgs)
		}
		if in.OIDCConfig.UsernameClaim != nil {
			config["username_claim"] = *in.OIDCConfig.UsernameClaim
		}
		if in.OIDCConfig.UsernamePrefix != nil {
			config["username_prefix"] = *in.OIDCConfig.UsernamePrefix
		}

		att["oidc_config"] = []interface{}{config}
	}

	return []interface{}{att}
}
