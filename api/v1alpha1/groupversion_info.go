package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var GroupVersion = schema.GroupVersion{
	Group:   "cloud.leo.io",
	Version: "v1alpha1",
}

var SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)

var AddToScheme = SchemeBuilder.AddToScheme

// Scheme register
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(
		GroupVersion,
		&SimpleApp{},
		&SimpleAppList{},
	)

	metav1.AddToGroupVersion(
		scheme,
		GroupVersion,
	)

	return nil
}
