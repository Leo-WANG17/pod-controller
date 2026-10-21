package controller

import (
	"testing"

	cloudv1alpha1 "github.com/Leo-WANG17/pod-controller/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestDesiredDeployment(t *testing.T) {
	// Create a test variable
	test := &cloudv1alpha1.SimpleApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "default",
			Labels: map[string]string{
				"env": "test_production",
			},
		},
		Spec: cloudv1alpha1.SimpleAppSpec{
			Replicas: 3,
			Image:    "nginx:latest",
		},
	}

	// Call func describedDeployment to create a deployment
	testDeployment := describedDeployment(test)

	// Assertion 1: check Name
	if testDeployment.Name != test.Name {
		t.Errorf(
			"Deployment name: got %q, expected %q",
			testDeployment.Name,
			test.Name,
		)
	}

	// Assertion 2: check Namespace
	if testDeployment.Namespace != test.Namespace {
		t.Errorf(
			"Deployment namespace: got %q, expected %q",
			testDeployment.Namespace,
			test.Namespace,
		)
	}

	// Assertion 3: check Replicas
	if testDeployment.Spec.Replicas == nil {
		t.Fatal("Deployment replicas must not be nil")
	}

	if *testDeployment.Spec.Replicas != 3 {
		t.Errorf("Deployment replicas: got %d, expected %d",
			*testDeployment.Spec.Replicas,
			3,
		)
	}

	// Assertion 4: check Selector Lable
	if testDeployment.Spec.Selector == nil {
		t.Fatal("Deployment selector must not be nil")
	}

	if got := testDeployment.Spec.Selector.MatchLabels["app"]; got != "test" {
		t.Errorf("Deployment selector: got %q, expected %q",
			got,
			"test",
		)
	}

	// Assertion 5: check Pod Template Label
	if got := testDeployment.Spec.Template.Labels["app"]; got != "test" {
		t.Errorf("Deployment template: got %q, expected %q",
			got,
			"test",
		)
	}

	// Assertion 6: check the number of Container
	if len(testDeployment.Spec.Template.Spec.Containers) != 1 {
		t.Fatal("Deployment template containers must have length 1")
	}

	// Assertion 7: check the image of Container
	if got := testDeployment.Spec.Template.Spec.Containers[0].Image; got != "nginx:latest" {
		t.Errorf("Container image: got %q, want %q",
			got,
			"nginx:latest",
		)
	}

	// Test: independency of ptr Replicas
	// Modify Deployment
	*testDeployment.Spec.Replicas = 5

	// Check Deployment
	if *testDeployment.Spec.Replicas != 5 {
		t.Errorf("Deployment replicas: got %d, expected %d",
			*testDeployment.Spec.Replicas,
			5,
		)
	}

	// Check original SimpleApp
	if test.Spec.Replicas != 3 {
		t.Errorf(
			"Deployment replicas: got %d, expected %d",
			test.Spec.Replicas,
			3,
		)
	}
}
