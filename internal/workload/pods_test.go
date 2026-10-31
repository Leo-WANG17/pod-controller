package workload

import (
	"context"
	"testing"

	cloudv1alpha1 "github.com/Leo-WANG17/pod-controller/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestListDeploymentPods(t *testing.T) {
	// Register Kubernetes Scheme
	scheme := runtime.NewScheme()

	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}

	if err := cloudv1alpha1.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}

	// Create target Deployment
	// Name: my-app
	// Namespace: default
	// UID: deployment-123
	// Selector: app=my-app
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-app",
			Namespace: "default",
			UID:       types.UID("deployment-123"),
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "my-app",
				},
			},
		},
	}

	// Create two owned ReplicaSets
	// rs-old and rs-new
	// Each must have its own UID
	// Controller OwnerReference -> deployment-123
	rsOld := &appsv1.ReplicaSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "rs-old",
			Namespace: "default",
			UID:       types.UID("rs-old-123"),
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "apps/v1",
					Kind:       "Deployment",
					Name:       "my-app",
					UID:        deployment.UID,
					Controller: ptr.To(true),
				},
			},
		},
	}

	rsNew := &appsv1.ReplicaSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "rs-new",
			Namespace: "default",
			UID:       types.UID("rs-new-123"),
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "apps/v1",
					Kind:       "Deployment",
					Name:       "my-app",
					UID:        deployment.UID,
					Controller: ptr.To(true),
				},
			},
		},
	}
	// Create two Pods
	// pod-old -> rs-old
	// pod-new -> rs-new
	// Both have label app=my-app
	podOld := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod-old",
			Namespace: "default",
			Labels: map[string]string{
				"app": "my-app",
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "apps/v1",
					Kind:       "ReplicaSet",
					Name:       "rs-old",
					UID:        rsOld.UID,
					Controller: ptr.To(true),
				},
			},
		},
	}

	podNew := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod-new",
			Namespace: "default",
			Labels: map[string]string{
				"app": "my-app",
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "apps/v1",
					Kind:       "ReplicaSet",
					Name:       "rs-new",
					UID:        rsNew.UID,
					Controller: ptr.To(true),
				},
			},
		},
	}

	// Create foreign ReplicaSet and Pod
	// Foreign Pod also has label app=my-app
	// But foreign RS belongs to another Deployment
	rsForeign := &appsv1.ReplicaSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "rs-foreign",
			Namespace: "default",
			UID:       types.UID("rs-foreign-123"),
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "apps/v1",
					Kind:       "Deployment",
					Name:       "my-app",
					UID:        types.UID("other-deployment-456"),
					Controller: ptr.To(true),
				},
			},
		},
	}

	podForeign := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod-foreign",
			Namespace: "default",
			Labels: map[string]string{
				"app": "my-app",
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "apps/v1",
					Kind:       "ReplicaSet",
					Name:       "rs-foreign",
					UID:        rsForeign.UID,
					Controller: ptr.To(true),
				},
			},
		},
	}

	// Build Fake Client with all objects
	fakeClient := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(
			deployment,
			rsOld,
			rsNew,
			rsForeign,
			podOld,
			podNew,
			podForeign,
		).
		Build()

	// Call ListDeploymentPods()
	pods, err := ListDeploymentPods(
		context.Background(),
		fakeClient,
		deployment,
	)
	if err != nil {
		t.Fatalf("ListDeploymentPods failed: %v", err)
	}

	// Assert returned Pods
	// Expected: pod-old and pod-new
	// Unexpected: foreign-pod
	gotNames := make(map[string]bool)

	for _, pod := range pods {
		gotNames[pod.Name] = true
	}

	if len(pods) != 2 {
		t.Errorf("Expected 2 Pods, got %d", len(pods))
	}

	if !gotNames["pod-old"] {
		t.Error("Expected pod-old to be returned")
	}

	if !gotNames["pod-new"] {
		t.Error("Expected pod-new to be returned")
	}

	if gotNames["foreign-pod"] {
		t.Error("Foreign Pod must not be returned")
	}
}
