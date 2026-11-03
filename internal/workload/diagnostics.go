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

		// Check whether containerStatus.State.Waiting is nil.
		if containerStatus.State.Waiting == nil {
			continue
		}

		// Get Waiting.Reason and Waiting.Message.
		reason := containerStatus.State.Waiting.Reason
		message := containerStatus.State.Waiting.Message

		// If Reason is ImagePullBackOff or CrashLoopBackOff,
		// append a PodDiagnosis.
		// Preserve PodName and ContainerName.
		switch reason {
		case "ImagePullBackOff":
			diagnoses = append(diagnoses, PodDiagnosis{
				PodName:       pod.Name,
				ContainerName: containerStatus.Name,
				Reason:        reason,
				Message:       message,
			})
		case "CrashLoopBackOff":
			diagnoses = append(diagnoses, PodDiagnosis{
				PodName:       pod.Name,
				ContainerName: containerStatus.Name,
				Reason:        reason,
				Message:       message,
			})
		}

	}

	return diagnoses
}
