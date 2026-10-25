package v1alpha1

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Unit Test 1 -- SimpleApp.DeepCopy()
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
			Image:    "nginx:latest",
		},
	}

	// Add conditions
	original.Status.Conditions = []metav1.Condition{
		{
			Type:    "Ready",
			Status:  metav1.ConditionTrue,
			Reason:  "AllReplicasReady",
			Message: "All replicas are ready",
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

	// test conditions
	// Modify the copied condition
	copied.Status.Conditions[0].Status = metav1.ConditionFalse

	// Verify the original condition remains unchanged
	if original.Status.Conditions[0].Status != metav1.ConditionTrue {
		t.Errorf(
			"Original condition changed unexpectedly: got %q, want %q",
			original.Status.Conditions[0].Status,
			metav1.ConditionTrue,
		)
	}
}

// Unit Test 2 -- SimpleAppList.DeepCopy()
func TestSimpleAppListDeepCopy(t *testing.T) {
	original := &SimpleAppList{
		Items: []SimpleApp{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "app-1",
					Labels: map[string]string{
						"env": "production",
					},
				},
				Spec: SimpleAppSpec{
					Replicas: 3,
					Image:    "nginx:latest",
				},
			},
		},
	}

	copied := original.DeepCopy()

	copied.Items[0].Name = "app-2"
	copied.Items[0].Labels["env"] = "development"
	copied.Items[0].Spec.Replicas = 5

	copied.Items = append(copied.Items, SimpleApp{
		ObjectMeta: metav1.ObjectMeta{
			Name: "app-3",
			Labels: map[string]string{
				"env": "testing",
			},
		},
		Spec: SimpleAppSpec{
			Replicas: 2,
			Image:    "redis:latest",
		},
	})

	// Confirm original.Items[0].Name = "app-1"
	if original.Items[0].Name != "app-1" {
		t.Errorf("original Name changed unexpectedly: %s", original.Items[0].Name)
	}
	// Confirm original.Items[0].Labels["env"] = "production"
	if original.Items[0].Labels["env"] != "production" {
		t.Errorf("original Labels changed unexpectedly: %v", original.Items[0].Labels)
	}
	// Confirm original.Items[0].Spec.Replicas = 3
	if original.Items[0].Spec.Replicas != 3 {
		t.Errorf("original Spec changed unexpectedly: %d", original.Items[0].Spec.Replicas)
	}
	// Confirm the number of element of original.Item will not be changed
	if len(original.Items) != 1 {
		t.Errorf("original Items length: got %d, want 1", len(original.Items))
	}

	if len(copied.Items) != 2 {
		t.Errorf("copied Items length: got %d, want 2", len(copied.Items))
	}
}

// Nil Receiver
func TestSimpleAppDeepCopyNil(t *testing.T) {
	var app *SimpleApp

	copied := app.DeepCopy()

	if copied != nil {
		t.Errorf("expected nil, got %v", copied)
	}
}
