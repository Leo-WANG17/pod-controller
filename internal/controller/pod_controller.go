package controller

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PodReconciler struct {
	Client client.Client
}

// Reconcile handles Pod changes and prints the CPU request.
func (r *PodReconciler) Reconcile(
	ctx context.Context,
	req ctrl.Request,
) (ctrl.Result, error) {
	var pod corev1.Pod

	// Handle the situation of "Pod Not Found"
	if err := r.Client.Get(ctx, req.NamespacedName, &pod); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	fmt.Println("Pod:", pod.Name)
	fmt.Println("Namespace:", pod.Namespace)

	if len(pod.Spec.Containers) == 0 {
		return ctrl.Result{}, nil
	}

	cpu, ok := pod.Spec.Containers[0].Resources.Requests[corev1.ResourceCPU]

	// Handle the situation of CPU is not set, then its value 0
	if !ok {
		fmt.Println("CPU request is not set")
		return ctrl.Result{}, nil
	}

	fmt.Println("CPU request:", cpu.String())

	return ctrl.Result{}, nil
}

func (r *PodReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		// For is in order to identify the Primary Resource of this controller
		// In this case, the Primary Resource is Pod
		// The working flow is when Pod changes -> enqueue -> Request -> Reconcile
		// Request is generally the Key of Primary Resource
		For(&corev1.Pod{}).
		Complete(r)
}
