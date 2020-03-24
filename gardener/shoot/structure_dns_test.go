package shoot

import (
	"testing"

	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	cmp "github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/schema"
)

func getDNSTestData() (map[string]interface{}, *corev1beta1.DNS) {
	domain := "foo.bar"
	dnsRaw := map[string]interface{}{
		"domain": domain,
	}
	dns := &corev1beta1.DNS{
		Domain: &domain,
	}

	return dnsRaw, dns
}

func TestExpandDNS(t *testing.T) {
	dnsRaw, dns := getDNSTestData()
	data := schema.TestResourceDataRaw(t, DNSResource().Schema, dnsRaw)
	out := ExpandDNS([]interface{}{data.Get("")})
	if diff := cmp.Diff(dns, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}

func TestFlattenDNS(t *testing.T) {
	dnsRaw, dns := getDNSTestData()
	out := FlattenDNS(dns)[0].(map[string]interface{})
	data := schema.TestResourceDataRaw(t, DNSResource().Schema, dnsRaw).Get("")
	if diff := cmp.Diff(data, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}
