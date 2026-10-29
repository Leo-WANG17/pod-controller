package controller

import (
	"context"
	"testing"

	cloudv1alpha1 "github.com/Leo-WANG17/pod-controller/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
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

// Test Deployment update reconciliation.
// Verify that Reconcile updates an existing Deployment
// when its replicas or image differ from the SimpleApp spec.
func TestReconcileUpdatesDeployment(t *testing.T) {
	// Initialize Scheme.
	scheme := runtime.NewScheme()

	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}

	if err := cloudv1alpha1.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	// Create a SimpleApp.
	// Expected replicas: 3
	// Expected image: nginx:latest
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

	// Create an existing Deployment.
	// Actual replicas: 5
	// Actual image: nginx:lastest
	deployment := describedDeployment(app)
	deployment.Spec.Replicas = ptr.To(int32(5))
	deployment.Spec.Template.Spec.Containers[0].Image = "redis:latest"

	// Set Controller OwnerReference.
	if err := controllerutil.SetControllerReference(
		app,
		deployment,
		scheme,
	); err != nil {
		t.Fatal(err)
	}

	// Initialize Fake Client with both objects.
	fakeClient := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(app, deployment).
		WithStatusSubresource(app, deployment).
		Build()

	// Initialize SimpleAppReconciler.
	reconcile := &SimpleAppReconciler{
		Client: fakeClient,
		Scheme: scheme,
	}

	// Call Reconcile().
	if _, err := reconcile.Reconcile(
		context.Background(),
		ctrl.Request{
			NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "test-app",
			},
		},
	); err != nil {
		t.Fatal("Reconcile failed")
	}

	// Retrieve the updated Deployment.
	var updatedDeployment appsv1.Deployment

	if err := fakeClient.Get(
		context.Background(),
		types.NamespacedName{
			Namespace: "default",
			Name:      "test-app",
		},
		&updatedDeployment,
	); err != nil {
		t.Fatal("failed to retrieve deployment")
	}

	// Assert replicas == 3 and image == nginx:latest.
	if updatedDeployment.Spec.Replicas == nil {
		t.Fatal("Deployment.Spec.Replicas is nil")
	}

	if *updatedDeployment.Spec.Replicas != 3 {
		t.Errorf(
			"Deployment replicas = %d, want 3",
			*updatedDeployment.Spec.Replicas,
		)
	}

	if len(updatedDeployment.Spec.Template.Spec.Containers) == 0 {
		t.Fatal("Deployment has no containers")
	}

	if updatedDeployment.Spec.Template.Spec.Containers[0].Image != "nginx:latest" {
		t.Errorf(
			"Deployment image = %s, want nginx:latest",
			updatedDeployment.Spec.Template.Spec.Containers[0].Image,
		)
	}

	// Verify OwnerReference remains correct.
	owner := metav1.GetControllerOf(&updatedDeployment)

	if owner == nil {
		t.Fatal("Deployment has no Controller OwnerReference")
	}

	if owner.UID != app.UID {
		t.Errorf(
			"Deployment owner UID = %s, want %s",
			owner.UID,
			app.UID,
		)
	}

	// Call Reconcile() again to test repeated reconciliation.
	resourceVersionBefore := updatedDeployment.ResourceVersion
	var deploymentAfter appsv1.Deployment

	if err := fakeClient.Get(
		context.Background(),
		types.NamespacedName{
			Namespace: "default",
			Name:      "test-app",
		},
		&deploymentAfter,
	); err != nil {
		t.Fatal(err)
	}

	if deploymentAfter.ResourceVersion != resourceVersionBefore {
		t.Errorf(
			"Deployment was unexpectedly updated during second Reconcile",
		)
	}
}

// Test ownership protection.
// Verify that Reconcile refuses to modify a Deployment
// controlled by a different SimpleApp.
func TestReconcileRejectsForeignDeployment(t *testing.T) {
	// Register Scheme
	scheme := runtime.NewScheme()

	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}

	if err := cloudv1alpha1.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}

	// Create SimpleApp.
	// Desired: 3 replicas, nginx:latest.
	app := &cloudv1alpha1.SimpleApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:       "test-app",
			Namespace:  "default",
			UID:        types.UID("simpleapp-uid"),
			Generation: 1,
		},
		Spec: cloudv1alpha1.SimpleAppSpec{
			Replicas: 3,
			Image:    "nginx:latest",
		},
	}

	// Create an existing Deployment.
	// Actual: 1 replica, redis:latest.
	foreignDeployment := describedDeployment(app)
	foreignDeployment.Spec.Replicas = ptr.To(int32(1))
	foreignDeployment.Spec.Template.Spec.Containers[0].Image = "redis:latest"

	// Set a FOREIGN Controller OwnerReference.
	foreignDeployment.OwnerReferences = []metav1.OwnerReference{
		{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
			Name:       "foreign-owner",
			UID:        "foreign-uid",
			Controller: ptr.To(true),
		},
	}

	// Create Fake Client with both objects.
	fakeClient := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(app, foreignDeployment).
		WithStatusSubresource(app, foreignDeployment).
		Build()

	// Create SimpleAppReconciler.
	reconcile := &SimpleAppReconciler{
		Client: fakeClient,
		Scheme: scheme,
	}

	// Call Reconcile().
	// Expect a non-nil error.
	if _, err := reconcile.Reconcile(
		context.Background(),
		ctrl.Request{
			NamespacedName: types.NamespacedName{
				Namespace: "default",
				Name:      "test-app",
			},
		},
	); err == nil {
		t.Fatal("Expected Reconcile to reject foreign Deployment")
	}

	// Get Deployment again from Fake Client.
	var deployment appsv1.Deployment

	if err := fakeClient.Get(
		context.Background(),
		types.NamespacedName{
			Namespace: "default",
			Name:      "test-app",
		},
		&deployment,
	); err != nil {
		t.Fatalf("Failed to get deployment: %v", err)
	}

	// Assert:
	// Replicas remains 1.
	if deployment.Spec.Replicas == nil {
		t.Fatal("Deployment replicas must not be nil")
	}

	if *deployment.Spec.Replicas != 1 {
		t.Errorf(
			"Foreign Deployment replicas changed: got %d, want 1",
			*deployment.Spec.Replicas,
		)
	}
	// Image remains redis:latest.
	if len(deployment.Spec.Template.Spec.Containers) == 0 {
		t.Fatal("Foreign Deployment template containers must not be empty")
	}

	if deployment.Spec.Template.Spec.Containers[0].Image != "redis:latest" {
		t.Errorf(
			"Foreign Deployment image changed: got %s, want redis:latest",
			deployment.Spec.Template.Spec.Containers[0].Image,
		)
	}
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

	// Case 3: Desired replicas exceed current replicas
	// Deployment status fields are all 3
	deployment.Status.UpdatedReplicas = 3
	// Call isDeploymentReady(deployment, 5)
	// Expected result: false
	verify = isDeploymentReady(deployment, 5)
	if verify {
		t.Errorf("Expected Deployment not to be ready")
	}

	// Table-Driven Tests
	tests := []struct {
		name            string
		deployment      appsv1.DeploymentStatus
		desiredReplicas int32
		want            bool
	}{
		{
			name: "AllReady",
			deployment: appsv1.DeploymentStatus{
				Replicas:          3,
				UpdatedReplicas:   3,
				ReadyReplicas:     3,
				AvailableReplicas: 3,
			},
			desiredReplicas: 3,
			want:            true,
		},
		{
			name: "UpdatedReplicas",
			deployment: appsv1.DeploymentStatus{
				Replicas:          3,
				UpdatedReplicas:   2,
				ReadyReplicas:     3,
				AvailableReplicas: 3,
			},
			desiredReplicas: 3,
			want:            false,
		},
		{
			name: "DesiredReplicas",
			deployment: appsv1.DeploymentStatus{
				Replicas:          3,
				UpdatedReplicas:   3,
				ReadyReplicas:     3,
				AvailableReplicas: 3,
			},
			desiredReplicas: 5,
			want:            false,
		},
		{
			name: "NotAvailable",
			deployment: appsv1.DeploymentStatus{
				Replicas:          3,
				UpdatedReplicas:   3,
				ReadyReplicas:     3,
				AvailableReplicas: 2,
			},
			desiredReplicas: 3,
			want:            false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deployment := &appsv1.Deployment{
				Status: tt.deployment,
			}

			got := isDeploymentReady(deployment, tt.desiredReplicas)

			if got != tt.want {
				t.Errorf(
					"isDeploymentReady() = %v, want %v",
					got,
					tt.want,
				)
			}
		})
	}
}

// Test case for complete Ready Condition
func TestReconcileReadyCondition(t *testing.T) {
	// Reconcile Status Integration-style Test
	tests := []struct {
		name                 string
		deploymentStatus     appsv1.DeploymentStatus
		wantCondition        metav1.ConditionStatus
		wantReason           string
		initialReadyReplicas int32
		wantReadyReplicas    int32
	}{
		// AllReady
		{
			name: "AllReady",
			deploymentStatus: appsv1.DeploymentStatus{
				ObservedGeneration: 2,
				Replicas:           3,
				UpdatedReplicas:    3,
				ReadyReplicas:      3,
				AvailableReplicas:  3,
			},
			wantCondition:        metav1.ConditionTrue,
			wantReason:           "AllReplicasReady",
			initialReadyReplicas: 0,
			wantReadyReplicas:    3,
		},
		// ReplicasNotReady
		{
			name: "ReplicasNotReady",
			deploymentStatus: appsv1.DeploymentStatus{
				ObservedGeneration: 2,
				Replicas:           3,
				UpdatedReplicas:    2,
				ReadyReplicas:      3,
				AvailableReplicas:  3,
			},
			wantCondition:        metav1.ConditionFalse,
			wantReason:           "ReplicasNotReady",
			initialReadyReplicas: 0,
			wantReadyReplicas:    3,
		},
		// DeploymentStatusStale
		{
			name: "DeploymentStatusStale",
			deploymentStatus: appsv1.DeploymentStatus{
				ObservedGeneration: 1,
				Replicas:           3,
				UpdatedReplicas:    3,
				ReadyReplicas:      3,
				AvailableReplicas:  3,
			},
			wantCondition:        metav1.ConditionFalse,
			wantReason:           "DeploymentStatusStale",
			initialReadyReplicas: 1,
			wantReadyReplicas:    1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare Scheme
			scheme := runtime.NewScheme()

			if err := clientgoscheme.AddToScheme(scheme); err != nil {
				t.Fatal(err)
			}

			if err := cloudv1alpha1.AddToScheme(scheme); err != nil {
				t.Fatal(err)
			}

			// Create a SimpleApp
			app := &cloudv1alpha1.SimpleApp{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "test-app",
					Namespace:  "default",
					UID:        types.UID("test-uid-123"),
					Generation: 1,
				},
				Spec: cloudv1alpha1.SimpleAppSpec{
					Replicas: 3,
					Image:    "nginx:latest",
				},
			}

			// Create an existing Deployment
			deployment := describedDeployment(app)

			// Set correct OwnerReference
			if err := controllerutil.SetControllerReference(
				app,
				deployment,
				scheme,
			); err != nil {
				t.Fatal(err)
			}

			// Set Deployment Status
			// All replicas ready
			deployment.Generation = 2

			deployment.Status = appsv1.DeploymentStatus{
				ObservedGeneration: 2,
				Replicas:           3,
				UpdatedReplicas:    3,
				ReadyReplicas:      3,
				AvailableReplicas:  3,
			}

			// Build Fake Client
			fakeClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(app, deployment).
				WithStatusSubresource(app, deployment).
				Build()

			// Ensure Deployment Spec matches SimpleApp Spec
			owner := metav1.GetControllerOf(deployment)

			if owner == nil {
				t.Fatal("Deployment has no controller OwnerReference")
			}

			if owner.UID != app.UID {
				t.Errorf(
					"Deployment owner UID = %s, want %s",
					owner.UID,
					app.UID,
				)
			}

			// Create SimpleAppReconciler
			reconciler := &SimpleAppReconciler{
				Client: fakeClient,
				Scheme: scheme,
			}

			// Call Reconcile()
			// Update SimpleApp
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

			if result != (ctrl.Result{}) {
				t.Errorf("Unexpected reconcile result: %+v", result)
			}

			// Get updated simpleapp
			var updatedApp cloudv1alpha1.SimpleApp

			if err := fakeClient.Get(
				context.Background(),
				types.NamespacedName{
					Name:      "test-app",
					Namespace: "default",
				},
				&updatedApp,
			); err != nil {
				t.Fatal(err)
			}

			// Assert Ready Condition is True
			condition := meta.FindStatusCondition(
				updatedApp.Status.Conditions,
				"Ready",
			)

			if condition == nil {
				t.Fatal("No ready condition")
			}

			if condition.Status != metav1.ConditionTrue {
				t.Errorf("Unexpected condition status: %v", condition.Status)
			}

			if condition.Reason != "AllReplicasReady" {
				t.Errorf("Unexpected condition reason: %v", condition.Reason)
			}

			if updatedApp.Status.ReadyReplicas != 3 {
				t.Errorf("Unexpected ready replicas: %v", updatedApp.Status.ReadyReplicas)
			}

			if updatedApp.Status.ObservedGeneration != updatedApp.Generation {
				t.Errorf("Unexpected observed generation: %v", updatedApp.Status.ObservedGeneration)
			}
		})
	}
}
