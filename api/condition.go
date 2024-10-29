package api

import (
	"fmt"
	"sort"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// A Conditioned may have conditions set or retrieved. Conditions
// indicate the status of a particular FSM state of a resource,
// or in the case of the condition of type "Ready", the overall
// status of the resource.
type Conditioned interface {
	// GetGeneration returns the `metadata.generation` of the Kubernetes resource on which these status conditions live.
	GetGeneration() int64
	// GetConditions returns the status conditions of the resource.
	GetConditions() []Condition
	// SetConditions sets the status conditions of the resource.
	SetConditions(c ...Condition)
	// GetCondition returns the status condition of the resource with the supplied type.
	GetCondition(ConditionType) Condition
}

// A ConditionType represents a condition a resource could be in.
type ConditionType string

// String returns ConditionType as a string
func (c ConditionType) String() string {
	return string(c)
}

// Condition types.
const (
	// TypeReady represents whether the resource has been successfully and completely processed.
	TypeReady ConditionType = "Ready"

	// TypeSynced resources are believed to be in sync with the
	// Kubernetes resources that manage their lifecycle.
	TypeSynced ConditionType = "Synced"

	// TypeReferencesValid indicates whether object references are valid (i.e. that they exist).
	TypeReferencesValid = "ReferencesValid"

	// ReasonReferencesExist is the reason that ReferencesValid is true.
	ReasonReferencesExist = "ReferencedObjectsExist"
)

// A ConditionReason represents the reason a resource is in a condition.
type ConditionReason string

// Reasons a resource is or is not ready.
const (
	ReasonAvailable   ConditionReason = "Available"
	ReasonUnavailable ConditionReason = "Unavailable"
	ReasonCreating    ConditionReason = "Creating"
	ReasonDeleting    ConditionReason = "Deleting"
)

// Reasons a resource is or is not synced.
const (
	ReasonReconcileSuccess ConditionReason = "ReconcileSuccess"
	ReasonReconcileError   ConditionReason = "ReconcileError"
)

// A Condition that may apply to a resource.
// +kubebuilder:object:generate=true
type Condition struct {
	// Type of this condition. At most one of each condition type may apply to
	// a resource at any point in time.
	Type ConditionType `json:"type"`

	// Status of this condition; is it currently True, False, or Unknown?
	Status corev1.ConditionStatus `json:"status"`

	// ObservedGeneration is the .metadata.generation that the condition was set based on.
	// For instance, if .metadata.generation is currently 12, but the
	// .status.conditions[x].observedGeneration is 9, the condition is out of date with respect
	// to the current state of the instance.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// LastTransitionTime is the last time this condition transitioned from one
	// status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`

	// A Reason for this condition's last transition from one status to another.
	Reason ConditionReason `json:"reason"`

	// A Message containing details about this condition's last transition from
	// one status to another, if any.
	// +optional
	Message string `json:"message,omitempty"`
}

// Equal returns true if the condition is identical to the supplied condition,
// ignoring the LastTransitionTime.
func (c Condition) Equal(other Condition) bool {
	return c.Type == other.Type &&
		c.Status == other.Status &&
		c.Reason == other.Reason &&
		c.Message == other.Message &&
		c.ObservedGeneration == other.ObservedGeneration
}

// WithMessage returns a condition by adding the provided message to existing
// condition.
func (c Condition) WithMessage(msg string) Condition {
	c.Message = msg
	return c
}

// IsEmpty returns true if the condition is empty.
func (c Condition) IsEmpty() bool {
	return c.Type == "" &&
		c.Status == "" &&
		c.Reason == "" &&
		c.Message == ""
}

// A ConditionedStatus reflects the observed status of a resource. Only
// one condition of each type may exist.
// +kubebuilder:object:generate=true
type ConditionedStatus struct {
	// Conditions of the resource.
	// +optional
	Conditions []Condition `json:"conditions,omitempty"`
}

// NewConditionedStatus returns a stat with the supplied conditions set.
func NewConditionedStatus(c ...Condition) *ConditionedStatus {
	s := &ConditionedStatus{}
	s.SetConditions(c...)
	return s
}

// GetConditions returns the condition for the given ConditionType if exists,
// otherwise returns nil
func (s *ConditionedStatus) GetConditions() []Condition {
	return s.Conditions
}

// GetCondition returns the condition for the given ConditionType if exists,
// otherwise returns nil
func (s *ConditionedStatus) GetCondition(ct ConditionType) Condition {
	for _, c := range s.Conditions {
		if c.Type == ct {
			return c
		}
	}

	return Condition{Type: ct, Status: corev1.ConditionUnknown}
}

// SetConditions sets the supplied conditions, replacing any existing conditions
// of the same type. This is a no-op if all supplied conditions are identical,
// ignoring the last transition time, to those already set.
// TODO(harveyxia) since this is invoked often for the fsm controller frame, improve efficiency by using hash map to make this O(len(c)) instead of O(len(c)*len(s.Conditions))
func (s *ConditionedStatus) SetConditions(c ...Condition) {
	for _, new := range c {
		exists := false
		for i, existing := range s.Conditions {
			if existing.Type != new.Type {
				continue
			}

			if existing.Equal(new) {
				exists = true
				continue
			}

			s.Conditions[i] = new
			exists = true
		}
		if !exists {
			s.Conditions = append(s.Conditions, new)
		}
	}
}

// Equal returns true if the status is identical to the supplied status,
// ignoring the LastTransitionTimes and order of statuses.
func (s *ConditionedStatus) Equal(other *ConditionedStatus) bool {
	if s == nil || other == nil {
		return s == nil && other == nil
	}

	if len(other.Conditions) != len(s.Conditions) {
		return false
	}

	sc := make([]Condition, len(s.Conditions))
	copy(sc, s.Conditions)

	oc := make([]Condition, len(other.Conditions))
	copy(oc, other.Conditions)

	// We should not have more than one condition of each type.
	sort.Slice(sc, func(i, j int) bool { return sc[i].Type < sc[j].Type })
	sort.Slice(oc, func(i, j int) bool { return oc[i].Type < oc[j].Type })

	for i := range sc {
		if !sc[i].Equal(oc[i]) {
			return false
		}
	}

	return true
}

// Creating returns a condition indicating the resource is currently
// being created.
func Creating() Condition {
	return Condition{
		Type:               TypeReady,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonCreating,
	}
}

// Deleting returns a condition indicating the resource is currently
// being deleted.
func Deleting() Condition {
	return Condition{
		Type:               TypeReady,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonDeleting,
	}
}

// Available returns a condition indicating the resource is
// currently observed to be available for use.
func Available() Condition {
	return Condition{
		Type:               TypeReady,
		Status:             corev1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonAvailable,
	}
}

// Unavailable returns a condition indicating the resource is not
// currently available for use. Unavailable should be set only when Crossplane
// expects the resource to be available but knows it is not, for example
// because its API reports it is unhealthy.
func Unavailable() Condition {
	return Condition{
		Type:               TypeReady,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonUnavailable,
	}
}

// ReconcileSuccess returns a condition indicating that Crossplane successfully
// completed the most recent reconciliation of the resource.
func ReconcileSuccess() Condition {
	return Condition{
		Type:               TypeSynced,
		Status:             corev1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonReconcileSuccess,
	}
}

// ReconcileError returns a condition indicating that Crossplane encountered an
// error while reconciling the resource. This could mean Crossplane was
// unable to update the resource to reflect its desired state, or that
// Crossplane was unable to determine the current actual state of the resource.
func ReconcileError(err error) Condition {
	return Condition{
		Type:               TypeSynced,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonReconcileError,
		Message:            err.Error(),
	}
}

// ReferencesValid returns a condition indicating that all object references
// are valid, i.e. that the referenced object exists.
func ReferencesValid() Condition {
	return Condition{
		Type:               TypeReferencesValid,
		LastTransitionTime: metav1.Now(),
		Status:             corev1.ConditionTrue,
		Reason:             ReasonReferencesExist,
		Message:            "All object references are valid.",
	}
}

// ReferencesInvalid returns a condition indicating that some object references
// are invalid, i.e. that they reference non-existent objects.
func ReferencesInvalid(reason ConditionReason, missingRefs []ObjectRef) Condition {
	var missingRefStrings []string
	for _, ref := range missingRefs {
		missingRefStrings = append(missingRefStrings, ref.ObjectKey().String())
	}

	return Condition{
		Type:               TypeReferencesValid,
		LastTransitionTime: metav1.Now(),
		Status:             corev1.ConditionFalse,
		Reason:             reason,
		Message:            fmt.Sprintf("Referenced objects are not found: %s", strings.Join(missingRefStrings, ", ")),
	}
}
