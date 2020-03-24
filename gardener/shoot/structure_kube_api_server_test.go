package shoot

import (
	"testing"

	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	cmp "github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/schema"
)

func getKubeAPIServerTestData() (map[string]interface{}, *corev1beta1.KubeAPIServerConfig) {
	basicAuth := true
	caBundle := "CABundle"
	clientID := "ClientID"
	groupsClaim := "GroupsClaim"
	groupsPrefix := "GroupsPrefix"
	issuerURL := "IssuerURL"
	usernameClaim := "UsernameClaim"
	usernamePrefix := "UsernamePrefix"
	kubeAPIServerRaw := map[string]interface{}{
		"enable_basic_authentication": basicAuth,
		"oidc_config": []interface{}{
			map[string]interface{}{
				"ca_bundle":       caBundle,
				"client_id":       clientID,
				"groups_claim":    groupsClaim,
				"groups_prefix":   groupsPrefix,
				"issuer_url":      issuerURL,
				"required_claims": map[string]interface{}{"key": "value"},
				"signing_algs":    []interface{}{"foo", "bar"},
				"username_claim":  usernameClaim,
				"username_prefix": usernamePrefix,
			},
		},
	}
	kubeAPIServer := &corev1beta1.KubeAPIServerConfig{
		EnableBasicAuthentication: &basicAuth,
		OIDCConfig: &corev1beta1.OIDCConfig{
			CABundle:       &caBundle,
			ClientID:       &clientID,
			GroupsClaim:    &groupsClaim,
			GroupsPrefix:   &groupsPrefix,
			IssuerURL:      &issuerURL,
			RequiredClaims: map[string]string{"key": "value"},
			SigningAlgs:    []string{"bar", "foo"},
			UsernameClaim:  &usernameClaim,
			UsernamePrefix: &usernamePrefix,
		},
	}

	return kubeAPIServerRaw, kubeAPIServer
}

func TestExpandKubeAPIServer(t *testing.T) {
	kubeAPIServerRaw, kubeAPIServer := getKubeAPIServerTestData()
	data := schema.TestResourceDataRaw(t, KubeAPIServerResource().Schema, kubeAPIServerRaw)
	out := ExpandKubeAPIServer([]interface{}{data.Get("")})
	if diff := cmp.Diff(kubeAPIServer, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}

func TestFlattenKubeAPIServer(t *testing.T) {
	kubeAPIServerRaw, kubeAPIServer := getKubeAPIServerTestData()
	out := FlattenKubeAPIServer(kubeAPIServer)[0].(map[string]interface{})
	data := schema.TestResourceDataRaw(t, KubeAPIServerResource().Schema, kubeAPIServerRaw).Get("")
	if diff := cmp.Diff(data, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}
