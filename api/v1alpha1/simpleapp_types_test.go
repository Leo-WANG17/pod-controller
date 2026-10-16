package v1alpha1

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSimpleAppDeepCopy(t *testing.T) {
	// Created a test object Simpleapp
	original := &SimpleApp{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-app",
			Labels: map[string]string{
				"env": "production",
			},
		},
		Spec: SimpleAppSpec{
			Replicas: 3,
			Image:    "nginx-latest",
		},
	}

	// Test operation
	copied := original.DeepCopy()
	copied.Name = "new-app"
	copied.Spec.Replicas = 5
	copied.Labels["env"] = "development"

	// Check 1: original.Name is still "my-app" or changed
	if original.Name != "my-app" {
		t.Errorf("original Name changed unexpectedly: %s", original.Name)
	}

	// Check 2: original.Lables["env"] is still "production" or changed
	if original.Labels["env"] != "production" {
		t.Errorf("original Labels changed unexpectedly: %v", original.Labels)
	}

	// Check 3: original.Spec.Replicas is still 3 or changed
	if original.Spec.Replicas != 3 {
		t.Errorf("original Spec changed unexpectedly: %d", original.Spec.Replicas)
	}
}
