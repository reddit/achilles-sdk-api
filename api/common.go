package api

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ClusterObjectRef references an object by name, namespace, and cluster.
// Used in multi-cluster APIs.
type ClusterObjectRef struct {
	// Name of the object. Required.
	Name string `json:"name"`

	// Namespace of the object. Required.
	Namespace string `json:"namespace"`

	// ClusterID of the object. Required.
	ClusterID string `json:"clusterId"`
}

// String returns the ClusterObjectRef as a string
func (o ClusterObjectRef) String() string {
	return strings.Join([]string{o.ClusterID, o.Namespace, o.Name}, string(types.Separator))
}

// ObjectRef references a namespace-scoped object by name and namespace.
type ObjectRef struct {
	// Name of the object. Required.
	Name string `json:"name"`

	// Namespace of the object. Required.
	Namespace string `json:"namespace"`
}

// ObjectKey returns the ObjectRef as a client.ObjectKey
func (o ObjectRef) ObjectKey() client.ObjectKey {
	return client.ObjectKey{Namespace: o.Namespace, Name: o.Name}
}

// ObjectRefFrom returns an *ObjectRef from a client.Object
func ObjectRefFrom(o client.Object) *ObjectRef {
	return &ObjectRef{
		Name:      o.GetName(),
		Namespace: o.GetNamespace(),
	}
}

// TypedObjectRef references an object by name and namespace and includes its Group, Version, and Kind.
type TypedObjectRef struct {

	// Group of the object. Required.
	Group string `json:"group"`

	// Version of the object. Required.
	Version string `json:"version"`

	// Kind of the object. Required.
	Kind string `json:"kind"`

	// Name of the object. Required.
	Name string `json:"name"`

	// Namespace of the object. Required.
	Namespace string `json:"namespace"`
}

func (t TypedObjectRef) GroupVersionKind() schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   t.Group,
		Version: t.Version,
		Kind:    t.Kind,
	}
}

func (t TypedObjectRef) ObjectKey() client.ObjectKey {
	return client.ObjectKey{
		Namespace: t.Namespace,
		Name:      t.Name,
	}
}

func (t TypedObjectRef) ObjectKeyNotSet() bool {
	return t.Name == "" && t.Namespace == ""
}

// ToCoreV1ObjectReference is a convenience method that returns a *corev1.ObjectReference with a subset of fields populated.
func (t TypedObjectRef) ToCoreV1ObjectReference() *corev1.ObjectReference {
	return &corev1.ObjectReference{
		Kind:      t.Kind,
		Name:      t.Name,
		Namespace: t.Namespace,
		APIVersion: v1.GroupVersion{
			Group:   t.Group,
			Version: t.Version,
		}.String(),
	}
}

func (t TypedObjectRef) String() string {
	return fmt.Sprintf("%s: %s", t.GroupVersionKind(), t.ObjectKey())
}

// NamedObjectRef references an object by name and optionally by namespace.
type NamedObjectRef struct {
	// Name of the object. Required.
	Name string `json:"name"`

	// Namespace of the object. Optional. Defaulting behavior is determined by the parent API.
	Namespace string `json:"namespace,omitempty"`
}
