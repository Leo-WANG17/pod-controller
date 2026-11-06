package workload

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestDiagnosePodMultipleContainers(t *testing.T) {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-app-test",
		},
		Status: corev1.PodStatus{
			ContainerStatuses: []corev1.ContainerStatus{
				{
					Name: "backend",
					State: corev1.ContainerState{
						Running: &corev1.ContainerStateRunning{},
					},
				},
				{
					Name: "worker",
					State: corev1.ContainerState{
						Waiting: &corev1.ContainerStateWaiting{
							Reason: "CrashLoopBackOff",
						},
					},
				},
				{
					Name: "sidecar",
					State: corev1.ContainerState{
						Waiting: &corev1.ContainerStateWaiting{
							Reason: "ImagePullBackOff",
						},
					},
				},
			},
		},
	}

	diagnoses := DiagnosePod(pod)

	// Verify that exactly 2 diagnoses were returned.
	if len(diagnoses) != 2 {
		t.Fatalf("expected 2 diagnoses, got %d", len(diagnoses))
	}

	foundWorker := false
	foundSidecar := false

	// Verify worker has CrashLoopBackOff.
	for _, diagnosis := range diagnoses {
		if diagnosis.PodName != "my-app-test" {
			t.Errorf(
				"expected PodName my-app-test, got %s",
				diagnosis.PodName,
			)
		}

		switch diagnosis.ContainerName {
		case "worker":
			foundWorker = true

			if diagnosis.Reason != "CrashLoopBackOff" {
				t.Errorf(
					"worker: expected CrashLoopBackOff, got %s",
					diagnosis.Reason,
				)
			}

		case "sidecar":
			foundSidecar = true

			if diagnosis.Reason != "ImagePullBackOff" {
				t.Errorf(
					"sidecar: expected ImagePullBackOff, got %s",
					diagnosis.Reason,
				)
			}

		default:
			t.Errorf(
				"unexpected diagnosis for container %s",
				diagnosis.ContainerName,
			)
		}
	}

	if !foundWorker {
		t.Error("expected diagnosis for worker")
	}

	if !foundSidecar {
		t.Error("expected diagnosis for sidecar")
	}

}
