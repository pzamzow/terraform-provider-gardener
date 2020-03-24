package metadata

import (
	"testing"

	cmp "github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getMetadataTestData() (map[string]interface{}, metav1.ObjectMeta) {
	metadataRaw := map[string]interface{}{
		"name":      "cluster-name",
		"namespace": "foo",
		"annotations": map[string]interface{}{
			"foo": "bar",
		},
		"labels": map[string]interface{}{
			"foo": "bar",
		},
		"resource_version": "1",
		"finalizers":       []interface{}{"foo", "bar"},
	}
	metadata := metav1.ObjectMeta{
		Name:      "cluster-name",
		Namespace: "foo",
		Annotations: map[string]string{
			"foo": "bar",
		},
		Labels: map[string]string{
			"foo": "bar",
		},
		ResourceVersion: "1",
		Finalizers:      []string{"bar", "foo"},
	}

	return metadataRaw, metadata
}

func TestExpandMetadata(t *testing.T) {
	metadataRaw, metadata := getMetadataTestData()
	data := schema.TestResourceDataRaw(t, MetadataResource("").Schema, metadataRaw)
	out := ExpandMetadata([]interface{}{data.Get("")})
	if diff := cmp.Diff(metadata, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}

func TestFlattenMetadata(t *testing.T) {
	metadataRaw, metadata := getMetadataTestData()
	out := FlattenMetadata(metadata)[0].(map[string]interface{})
	data := schema.TestResourceDataRaw(t, MetadataResource("").Schema, metadataRaw)
	if diff := cmp.Diff(data.Get(""), out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}
