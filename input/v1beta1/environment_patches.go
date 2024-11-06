package v1beta1

import (
	"github.com/crossplane/function-sdk-go/resource/composite"
	"github.com/fernandezcuesta/function-patch-and-transform/input/v1beta1"
	"github.com/fernandezcuesta/function-patch-and-transform/pt"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// EnvironmentPatch objects are applied between the composite resource and
// the environment. Their behaviour depends on the Type selected. The default
// Type, FromCompositeFieldPath, copies a value from the composite resource
// to the environment, applying any defined transformers.
type EnvironmentPatch struct {
	// Type sets the patching behaviour to be used. Each patch type may require
	// its own fields to be set on the Patch object.
	// +optional
	// +kubebuilder:validation:Enum=FromCompositeFieldPath;CombineFromComposite
	// +kubebuilder:default=FromCompositeFieldPath
	Type v1beta1.PatchType `json:"type,omitempty"`

	v1beta1.Patch `json:",inline"`
}

// GetType returns the patch type. If the type is not set, it returns the default type.
func (ep *EnvironmentPatch) GetType() v1beta1.PatchType {
	if ep.Type == "" {
		return v1beta1.PatchTypeFromCompositeFieldPath
	}
	return ep.Type
}

// ApplyEnvironmentPatch applies a patch to or from the environment. Patches to
// the environment are always from the observed XR. Patches from the environment
// are always to the desired XR.
func ApplyEnvironmentPatch(p *EnvironmentPatch, env *unstructured.Unstructured, oxr, dxr *composite.Unstructured) error {
	switch p.GetType() {
	// From observed XR to environment.
	case v1beta1.PatchTypeFromCompositeFieldPath,
		v1beta1.PatchTypeToEnvironmentFieldPath:
		return pt.ApplyFromFieldPathPatch(p, oxr, env)
	case v1beta1.PatchTypeCombineFromComposite:
		return pt.ApplyCombineFromVariablesPatch(p, oxr, env)

	// From environment to desired XR.
	case v1beta1.PatchTypeToCompositeFieldPath,
		v1beta1.PatchTypeFromEnvironmentFieldPath:
		return pt.ApplyFromFieldPathPatch(p, env, dxr)
	case v1beta1.PatchTypeCombineToComposite:
		return pt.ApplyCombineFromVariablesPatch(p, env, dxr)

	// Invalid patch types in this context.
	case v1beta1.PatchTypeCombineFromEnvironment,
		v1beta1.PatchTypeCombineToEnvironment:
		// Nothing to do.

	case v1beta1.PatchTypePatchSet:
		// Already resolved - nothing to do.
	}
	return nil
}
