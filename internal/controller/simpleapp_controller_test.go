package controller

import (
	"context"
	"testing"

	cloudv1alpha1 "github.com/Leo-WANG17/pod-controller/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
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

// Fake client test
func TestReconcileCreateDeployment(t *testing.T) {
	// Prepare test Scheme
	scheme := runtime.NewScheme()

	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}

	if err := cloudv1alpha1.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}

	// Fake Client create
	app := &cloudv1alpha1.SimpleApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-app",
			Namespace: "default",
			UID:       types.UID("test-uid-123"),
		},
		Spec: cloudv1alpha1.SimpleAppSpec{
			Replicas: 3,
			Image:    "nginx:latest",
		},
	}

	// This fake client without Deployment
	fakeClient := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(app).
		Build()

	// Creat Reconciler
	reconciler := &SimpleAppReconciler{
		Client: fakeClient,
		Scheme: scheme,
	}

	// Call Reconcile
	result, err := reconciler.Reconcile(
		context.Background(),
		ctrl.Request{
			NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "test-app",
			},
		},
	)
	if err != nil {
		t.Fatalf("Reconcile failed: %v", err)
	}

	_ = result

	if result != (ctrl.Result{}) {
		t.Errorf("Unexpected reconcile result: %+v", result)
	}

	// get deployment from fake client
	var deployment appsv1.Deployment

	err = fakeClient.Get(
		context.Background(),
		types.NamespacedName{
			Namespace: "default",
			Name:      "test-app",
		},
		&deployment,
	)

	if err != nil {
		t.Fatalf("Failed to get deployment: %v", err)
	}

	// Assertion 1
	if deployment.Name != "test-app" {
		t.Errorf(
			"Deployment name: got %q, expected %q",
			deployment.Name,
			"test-app",
		)
	}

	// Assertion 2
	replicas := deployment.Spec.Replicas

	if replicas == nil {
		t.Fatal("Deployment replicas must not be nil")
	}

	if *replicas != 3 {
		t.Errorf(
			"Deployment replicas: got %d, expected %d",
			*replicas,
			3,
		)
	}

	// Assertion 3
	containers := deployment.Spec.Template.Spec.Containers

	if len(containers) != 1 {
		t.Fatalf("Container count: got %d, want 1", len(containers))
	}

	if containers[0].Image != "nginx:latest" {
		t.Errorf(
			"Deployment template container: got %q, expected %q",
			containers[0].Image,
			"nginx:latest",
		)
	}

	// OwnerReference verification
	owner := metav1.GetControllerOf(&deployment)

	if owner == nil {
		t.Fatal("Deployment owner must not be nil")
	}

	if owner.Kind != "SimpleApp" {
		t.Errorf("Owner Kind: got %q, want %q", owner.Kind, "SimpleApp")
	}

	if owner.UID != app.UID {
		t.Errorf("Owner UID: got %q, want %q", owner.UID, app.UID)
	}

	if owner.APIVersion != cloudv1alpha1.GroupVersion.String() {
		t.Errorf(
			"Owner APIVersion: got %q, want %q",
			owner.APIVersion,
			cloudv1alpha1.GroupVersion.String(),
		)
	}
}

// TODO: Test Deployment update reconciliation.
// Verify that Reconcile updates an existing Deployment
// when its replicas or image differ from the SimpleApp spec.
func TestReconcileUpdatesDeployment(t *testing.T) {
	// TODO 1: Initialize Scheme.

	// TODO 2: Create a SimpleApp.
	// Expected replicas: 3
	// Expected image: nginx:latest

	// TODO 3: Create an existing Deployment.
	// Actual replicas: 5
	// Actual image: nginx:lastest

	// TODO 4: Set Controller OwnerReference.

	// TODO 5: Initialize Fake Client with both objects.

	// TODO 6: Initialize SimpleAppReconciler.

	// TODO 7: Call Reconcile().

	// TODO 8: Retrieve the updated Deployment.

	// TODO 9: Assert replicas == 3 and image == nginx:latest.

	// TODO 10: Verify OwnerReference remains correct.

	// TODO 11: Call Reconcile() again to test repeated reconciliation.
}

// TODO: Test ownership protection.
// Verify that Reconcile refuses to modify a Deployment
// controlled by a different SimpleApp.
func TestReconcileRejectsForeignDeployment(t *testing.T) {
	// TODO 1: Initialize Scheme.

	// TODO 2: Create a SimpleApp with UID "owner-a".

	// TODO 3: Create an existing Deployment.
	// Give its controller OwnerReference a different UID: "owner-b".

	// TODO 4: Initialize Fake Client with both objects.

	// TODO 5: Initialize SimpleAppReconciler.

	// TODO 6: Call Reconcile().

	// TODO 7: Assert that Reconcile returns an ownership error.

	// TODO 8: Retrieve Deployment again.

	// TODO 9: Verify that its spec was not modified.
}

// test Ready Condition
func isDeploymentReady(
	deployment *appsv1.Deployment,
	desiredReplicas int32,
) bool {
	return deployment.Status.ReadyReplicas == desiredReplicas &&
		deployment.Status.Replicas == desiredReplicas &&
		deployment.Status.AvailableReplicas == desiredReplicas &&
		deployment.Status.UpdatedReplicas == desiredReplicas
}

func TestIsDeploymentReady(t *testing.T) {
	// Expected replicas: 3
	// Set all Status replica fields to 3
	deployment := &appsv1.Deployment{
		Status: appsv1.DeploymentStatus{
			Replicas:          3,
			UpdatedReplicas:   3,
			ReadyReplicas:     3,
			AvailableReplicas: 3,
		},
	}

	// Call isDeploymentReady()
	// Expect true
	verify := isDeploymentReady(deployment, 3)

	if !verify {
		t.Errorf("Deployment is not ready")
	}
	// Change UpdatedReplicas to 2
	deployment.Status.UpdatedReplicas = 2

	// Call isDeploymentReady() again
	// Expect false
	verify = isDeploymentReady(deployment, 3)
	if verify {
		t.Errorf("Expected Deployment not to be ready")
	}

}
