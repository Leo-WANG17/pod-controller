package controller

import (
	"context"

	cloudv1alpha1 "github.com/Leo-WANG17/pod-controller/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type SimpleAppReconciler struct {
	Client client.Client
	Scheme *runtime.Scheme
}

// Reconcile the reconciler of SimpleApp
func (r *SimpleAppReconciler) Reconcile(
	ctx context.Context,
	req ctrl.Request,
) (ctrl.Result, error) {
	var app cloudv1alpha1.SimpleApp

	logger := ctrl.LoggerFrom(ctx)

	// Retrieve the SimpleApp from the controller-runtime cache.
	// Success then print the information
	if err := r.Client.Get(ctx, req.NamespacedName, &app); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	logger.V(1).Info(
		"Reconciling SimpleApp",
		"name", app.Name,
		"namespace", app.Namespace,
		"replicas", app.Spec.Replicas,
		"image", app.Spec.Image,
	)

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

	return ctrl.Result{}, nil
}

// describedDeployment Deployment constructor
func describedDeployment(app *cloudv1alpha1.SimpleApp) *appsv1.Deployment {
	replicas := app.Spec.Replicas

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": app.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": app.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "app",
							Image: app.Spec.Image,
						},
					},
				},
			},
		},
	}

	return deployment
}

func (r *SimpleAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cloudv1alpha1.SimpleApp{}).
		Complete(r)
}
