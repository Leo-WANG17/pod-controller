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
		t.Fatalf(
			"expected 2 diagnoses, got %d",
			len(diagnoses),
		)
	}

	// Verify worker has CrashLoopBackOff.
	for _, diagnosis := range diagnoses {
		if diagnosis.ContainerName == "worker" {
			if diagnosis.Reason != "CrashLoopBackOff" {
				t.Errorf("expected reason is CrashLoopBackOff, got %s", diagnosis.Reason)
			}
			// Verify sidecar has ImagePullBackOff.
		} else if diagnosis.ContainerName == "sidecar" {
			if diagnosis.Reason != "ImagePullBackOff" {
				t.Errorf("expected reason is ImagePullBackOff, got %s", diagnosis.Reason)
			}
			// Verify backend has no diagnosis.
		} else if diagnosis.ContainerName == "backend" {
			t.Errorf("%s has no diagnosis", diagnosis.ContainerName)
		}
	}

}
