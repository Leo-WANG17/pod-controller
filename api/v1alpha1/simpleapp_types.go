package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// SimpleAppSpec Serialization: Go filed -> JSON field
type SimpleAppSpec struct {
	Replicas int32  `json:"replicas"`
	Image    string `json:"image"`
}

// SimpleApp Singular resource
type SimpleApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec SimpleAppSpec `json:"spec,omitempty"`
}

func (s *SimpleApp) DeepCopyObject() runtime.Object {
	if s == nil {
		return nil
	}

	out := new(SimpleApp)
	*out = *s
	s.ObjectMeta.DeepCopyInto(&out.ObjectMeta)

	return out
}

func (in *SimpleApp) DeepCopyInto(out *SimpleApp) {
	*out = *in
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
}

// SimpleAppList Multiple resources of SimpleApp
type SimpleAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []SimpleApp `json:"items"`
}

// DeepCopyInto To handle every element of slice SimpleApp into SimpleAppList slice
func (in *SimpleAppList) DeepCopyInto(out *SimpleAppList) {
	*out = *in

	if in.Items != nil {
		out.Items = make([]SimpleApp, len(in.Items))

		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}
}

// DeepCopy Handle the deep copy process, by calling the func
func (in *SimpleAppList) DeepCopy() *SimpleAppList {
	if in == nil {
		return nil
	}

	out := new(SimpleAppList)
	in.DeepCopyInto(out)

	return out
}
