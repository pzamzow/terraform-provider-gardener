package metadata

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	lib "github.com/kyma-incubator/terraform-provider-gardener/gardener/lib"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MetadataResource(objectName string) *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"annotations": {
				Type:        schema.TypeMap,
				Description: fmt.Sprintf("An unstructured key value map stored with the %s that may be used to store arbitrary metadata. More info: http://kubernetes.io/docs/user-guide/annotations", objectName),
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ValidateFunc:     lib.ValidateAnnotations,
				DiffSuppressFunc: lib.SuppressGardenerAnnotations,
			},
			"generation": {
				Type:        schema.TypeInt,
				Description: "A sequence number representing a specific generation of the desired state.",
				Computed:    true,
			},
			"labels": {
				Type:             schema.TypeMap,
				Description:      fmt.Sprintf("Map of string keys and values that can be used to organize and categorize (scope and select) the %s. May match selectors of replication controllers and services. More info: http://kubernetes.io/docs/user-guide/labels", objectName),
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				ValidateFunc:     lib.ValidateLabels,
				DiffSuppressFunc: lib.SuppressGardenerLabels,
			},
			"name": {
				Type:        schema.TypeString,
				Description: fmt.Sprintf("Name of the %s, must be unique. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/identifiers#names", objectName),
				Optional:    true,

				Computed:     true,
				ValidateFunc: lib.ValidateName,
			},
			"resource_version": {
				Type:        schema.TypeString,
				Description: fmt.Sprintf("An opaque value that represents the internal version of this %s that can be used by clients to determine when %s has changed. Read more: https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#concurrency-control-and-consistency", objectName, objectName),
				Computed:    true,
			},
			"finalizers": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"self_link": {
				Type:        schema.TypeString,
				Description: fmt.Sprintf("A URL representing this %s.", objectName),
				Computed:    true,
			},
			"uid": {
				Type:        schema.TypeString,
				Description: fmt.Sprintf("The unique in time and space value for this %s. More info: http://kubernetes.io/docs/user-guide/identifiers#uids", objectName),
				Computed:    true,
			},
			"namespace": {
				Type:        schema.TypeString,
				Description: fmt.Sprintf("Namespace defines the space within which name of the %s must be unique.", objectName),
				Optional:    true,
				Default:     "default",
			},
		},
	}
}

func ExpandMetadata(in []interface{}) metav1.ObjectMeta {
	meta := metav1.ObjectMeta{}
	if len(in) < 1 {
		return meta
	}
	m := in[0].(map[string]interface{})

	if v, ok := m["annotations"].(map[string]interface{}); ok && len(v) > 0 {
		meta.Annotations = lib.ExpandStringMap(m["annotations"].(map[string]interface{}))
	}
	if v, ok := m["labels"].(map[string]interface{}); ok && len(v) > 0 {
		meta.Labels = lib.ExpandStringMap(m["labels"].(map[string]interface{}))
	}
	if v, ok := m["generate_name"].(string); ok {
		meta.GenerateName = v
	}
	if v, ok := m["name"].(string); ok {
		meta.Name = v
	}
	if v, ok := m["namespace"].(string); ok {
		meta.Namespace = v
	}
	if v, ok := m["resource_version"].(string); ok {
		meta.ResourceVersion = v
	}
	if v, ok := m["resource_version"].(string); ok {
		meta.ResourceVersion = v
	}
	if v, ok := m["finalizers"].(*schema.Set); ok {
		meta.Finalizers = lib.ExpandSet(v)
	}

	return meta
}

func FlattenMetadata(meta metav1.ObjectMeta) []interface{} {
	m := make(map[string]interface{})

	m["annotations"] = lib.FlattenStringMap(meta.Annotations)
	if meta.GenerateName != "" {
		m["generate_name"] = meta.GenerateName
	}
	m["labels"] = lib.FlattenStringMap(meta.Labels)
	m["name"] = meta.Name
	m["resource_version"] = meta.ResourceVersion
	m["finalizers"] = lib.NewStringSet(schema.HashString, meta.Finalizers)
	m["self_link"] = meta.SelfLink
	m["uid"] = fmt.Sprintf("%v", meta.UID)
	m["generation"] = int(meta.Generation)

	if meta.Namespace != "" {
		m["namespace"] = meta.Namespace
	}

	return []interface{}{m}
}
