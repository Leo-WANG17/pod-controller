package workload

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ListNamespacePods Get the list of ALL Pods within an indicated Namespace
func ListNamespacePods(
	ctx context.Context,
	k8sClient client.Client,
	namespace string,
) ([]corev1.Pod, error) {
	var podList corev1.PodList

	err := k8sClient.List(
		ctx,
		&podList,
		client.InNamespace(namespace),
	)

	if err != nil {
		return nil, err
	}

	return podList.Items, nil
}

// ListDeploymentPods Selecting the list of ALL pods which belong to an indicated Deployment
func ListDeploymentPods(
	ctx context.Context,
	k8sClient client.Client,
	deployment *appsv1.Deployment,
) ([]corev1.Pod, error) {
	// Validate deployment and deployment.Spec.Selector.
	if deployment == nil {
		return nil, fmt.Errorf("deployment can not be nil")
	}

	if deployment.Spec.Selector == nil {
		return nil, fmt.Errorf("deployment %s/%s has nil selector",
			deployment.Namespace,
			deployment.Name,
		)
	}

	if deployment.UID == "" {
		return nil, fmt.Errorf("deployment %s/%s has empty UID",
			deployment.Namespace,
			deployment.Name,
		)
	}

	// Convert LabelSelector using metav1.LabelSelectorAsSelector.
	selector, err := metav1.LabelSelectorAsSelector(
		deployment.Spec.Selector,
	)
	if err != nil {
		return nil, err
	}

	// List candidate Pods in deployment.Namespace.
	var podList corev1.PodList

	err = k8sClient.List(
		ctx,
		&podList,
		client.InNamespace(deployment.Namespace),
		client.MatchingLabelsSelector{
			Selector: selector,
		},
	)

	if err != nil {
		return nil, err
	}
	// List ReplicaSets in deployment.Namespace.
	var rsList appsv1.ReplicaSetList

	err = k8sClient.List(
		ctx,
		&rsList,
		client.InNamespace(deployment.Namespace),
	)

	if err != nil {
		return nil, err
	}

	// Build a set of ReplicaSet UIDs whose Controller Owner
	// points to this Deployment.
	ownerReplicaSets := make(map[types.UID]struct{})

	for i := range rsList.Items {
		rs := &rsList.Items[i]

		owner := metav1.GetControllerOf(rs)

		if owner == nil {
			continue
		}

		if owner.UID == deployment.UID &&
			owner.Kind == "Deployment" &&
			owner.APIVersion == appsv1.SchemeGroupVersion.String() {
			ownerReplicaSets[rs.UID] = struct{}{}
		}
	}

	// Iterate candidate Pods and verify their Controller Owner
	// points to one of the owned ReplicaSets.
	var verifiedPods []corev1.Pod

	for i := range podList.Items {
		pod := &podList.Items[i]

		owner := metav1.GetControllerOf(pod)

		if owner == nil {
			continue
		}

		if owner.Kind != "ReplicaSet" ||
			owner.APIVersion != appsv1.SchemeGroupVersion.String() {
			continue
		}

		_, exists := ownerReplicaSets[owner.UID]

		if exists {
			verifiedPods = append(verifiedPods, *pod)
		}
	}

	// Return the verified Pods.
	return verifiedPods, nil
}
