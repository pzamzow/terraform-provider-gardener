package shoot

import (
	"encoding/json"
	"testing"

	awsAlpha1 "github.com/gardener/gardener-extension-provider-aws/pkg/apis/aws/v1alpha1"
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	cmp "github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/schema"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func getProviderAWSInfrastructureConfigTestData() (map[string]interface{}, *corev1beta1.ProviderConfig) {
	vpcCIDR := "10.250.0.0/16"
	EnableECRAccess := true
	vpcID := "foo"
	awsConfig, _ := json.Marshal(awsAlpha1.InfrastructureConfig{
		TypeMeta: v1.TypeMeta{
			APIVersion: "aws.provider.extensions.gardener.cloud/v1alpha1",
			Kind:       "InfrastructureConfig",
		},
		EnableECRAccess: &EnableECRAccess,
		Networks: awsAlpha1.Networks{
			VPC: awsAlpha1.VPC{
				CIDR: &vpcCIDR,
				ID:   &vpcID,
			},
			Zones: []awsAlpha1.Zone{
				awsAlpha1.Zone{
					Name:     "eu-central-1a",
					Internal: vpcCIDR,
					Public:   vpcCIDR,
					Workers:  vpcCIDR,
				},
			},
		},
	})
	awsRaw := map[string]interface{}{
		"enableecraccess": true,
		"networks": []interface{}{
			map[string]interface{}{
				"vpc": []interface{}{
					map[string]interface{}{
						"cidr": vpcCIDR,
						"id":   vpcID,
					},
				},
				"zones": []interface{}{
					map[string]interface{}{
						"name":     "eu-central-1a",
						"internal": vpcCIDR,
						"public":   vpcCIDR,
						"workers":  vpcCIDR,
					},
				},
			},
		},
	}
	aws := &corev1beta1.ProviderConfig{
		RawExtension: runtime.RawExtension{
			Raw: awsConfig,
		},
	}

	return awsRaw, aws
}

func TestExpandProviderAWS(t *testing.T) {
	awsRaw, aws := getProviderAWSInfrastructureConfigTestData()
	data := schema.TestResourceDataRaw(t, AWSInfrastructureConfigResource().Schema, awsRaw)
	out := ExpandAWSInfrastructureConfig([]interface{}{data.Get("")})
	if diff := cmp.Diff(aws, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}

func TestFlattenProviderAWS(t *testing.T) {
	awsRaw, aws := getProviderAWSInfrastructureConfigTestData()
	out := FlattenAWSInfrastructureConfig(aws)[0].(map[string]interface{})
	data := schema.TestResourceDataRaw(t, AWSInfrastructureConfigResource().Schema, awsRaw).Get("")
	if diff := cmp.Diff(data, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}
