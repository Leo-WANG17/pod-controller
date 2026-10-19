package controller

import (
	"context"
	"fmt"

	cloudv1alpha1 "github.com/Leo-WANG17/pod-controller/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type SimpleAppReconciler struct {
	Client client.Client
}

func (r *SimpleAppReconciler) Reconcile(
	ctx context.Context,
	req ctrl.Request,
) (ctrl.Result, error) {
	var app cloudv1alpha1.SimpleApp

	logger := ctrl.LoggerFrom(ctx)

	logger.Info(
		"Reconciling SimpleApp",
		"name", app.Name,
		"namespace", app.Namespace,
		"replicas", app.Spec.Replicas,
		"image", app.Spec.Image,
	)

	// Retrieve the SimpleApp from the controller-runtime cache.
	// Success then print the information
	if err := r.Client.Get(ctx, req.NamespacedName, &app); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Handle miss required field (Replicas and Image) in this case
	if app.Spec.Replicas < 1 {
		logger.Info(
			"Invalid SimpleApp: replicas must be at least 1",
			"replicas", app.Spec.Replicas,
		)
		return ctrl.Result{}, nil
	}

	if app.Spec.Image == "" {
		logger.Info(
			"Invalid SimpleApp: image must not be empty",
		)
		return ctrl.Result{}, nil
	}

	fmt.Println("Replicas: ", app.Spec.Replicas)
	fmt.Println("Image: ", app.Spec.Image)

	return ctrl.Result{}, nil
}

func (r *SimpleAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cloudv1alpha1.SimpleApp{}).
		Complete(r)
}
