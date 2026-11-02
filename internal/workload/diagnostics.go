package workload

import corev1 "k8s.io/api/core/v1"

type PodDiagnosis struct {
	PodName       string
	ContainerName string
	Reason        string
	Message       string
}

func DiagnosePod(pod *corev1.Pod) []PodDiagnosis {
	if pod == nil {
		return nil
	}

	var diagnoses []PodDiagnosis

	for _, containerStatus := range pod.Status.ContainerStatuses {

		// TODO 1:
		// Check whether containerStatus.State.Waiting is nil.

		// TODO 2:
		// Get Waiting.Reason and Waiting.Message.

		// TODO 3:
		// If Reason is ImagePullBackOff or CrashLoopBackOff,
		// append a PodDiagnosis.

		// TODO 4:
		// Preserve PodName and ContainerName.
	}

	return diagnoses
}
