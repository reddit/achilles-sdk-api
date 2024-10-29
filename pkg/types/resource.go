package types

import (
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/reddit/achilles-sdk-api/api"
)

// FSMResource constrains to the types necessary for controller management.
type FSMResource[T any] interface {
	client.Object   // must be a k8s resource
	api.Conditioned // must have Conditions
	ResourceManager // must manage a set of child resources
	*T              // must be a pointer
}

// Resource constrains to the types necessary for controller management.
type Resource[T any] interface {
	client.Object   // must be a k8s resource
	api.Conditioned // must have Conditions
	*T              // must be a pointer
}

// ResourceManager is a k8s resource that manages a set of child resources.
type ResourceManager interface {
	// SetManagedResources sets the refs for child resources managed by the controller.
	SetManagedResources(refs []api.TypedObjectRef)
	// GetManagedResources gets the refs for child resources managed by the controller.
	GetManagedResources() []api.TypedObjectRef
}

// ClaimedResource is a k8s resource that can act as a Claimed.
type ClaimedResource interface {
	// GetClaimRef returns a reference to the claim that created this resource.
	GetClaimRef() *api.TypedObjectRef
	// SetClaimRef sets the reference to the claim that created this resource.
	SetClaimRef(ref *api.TypedObjectRef)
}

// ClaimResource is a k8s resource that can act as a Claim.
type ClaimResource interface {
	// GetClaimedRef returns a reference to the resource claimed by this claim.
	GetClaimedRef() *api.TypedObjectRef
	// SetClaimedRef sets the reference to the resource claimed by this claim.
	SetClaimedRef(ref *api.TypedObjectRef)
}

// ClaimedType constrains a resource to a ClaimedResource.
type ClaimedType[T any] interface {
	FSMResource[T]
	ClaimedResource
}

// ClaimType constrains a resource to a ClaimResource.
type ClaimType[T any] interface {
	Resource[T]
	ClaimResource
}
