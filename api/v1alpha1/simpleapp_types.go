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

func (s SimpleApp) DeepCopyObject() runtime.Object {
	//TODO implement me
	panic("implement me")
}

// SimpleAppList Multiple resources of SimpleApp
type SimpleAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []SimpleApp `json:"items"`
}

func (s SimpleAppList) DeepCopyObject() runtime.Object {
	//TODO implement me
	panic("implement me")
}
