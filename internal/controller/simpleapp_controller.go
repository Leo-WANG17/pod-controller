package controller

import (
	"context"
	"fmt"

	cloudv1alpha1 "github.com/Leo-WANG17/pod-controller/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
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

	// Handle deployment
	// To check if there is deployment existed, if no, create it
	// 1. Declare an existing Deployment
	var existing appsv1.Deployment

	// 2. Get Deployment
	err := r.Client.Get(
		ctx,
		// client.ObjectKey is used for identify the namespace and name of deployment
		client.ObjectKey{
			Namespace: app.Namespace,
			Name:      app.Name,
		},
		&existing,
	)

	// 3. NotFound -> SetControllerReference -> Create
	if apierrors.IsNotFound(err) {
		//     a. Call desiredDeployment(&app)
		desired := describedDeployment(&app)

		//     b. SetControllerReference(...)
		if err := controllerutil.SetControllerReference(
			&app,
			desired,
			r.Scheme,
		); err != nil {
			return ctrl.Result{}, err
		}

		//     c. Client.Create(...)
		if err := r.Client.Create(ctx, desired); err != nil {
			return ctrl.Result{},
				fmt.Errorf("failed to create deployment %s/%s: %w",
					desired.Namespace,
					desired.Name,
					err,
				)
		}

		logger.Info(
			"Deployment created",
			"namespace", desired.Namespace,
			"name", desired.Name,
		)

		return ctrl.Result{}, nil
	}

	// 4. If another error:
	if err != nil {
		return ctrl.Result{}, err
	}

	// 5. If Deployment already exists here.
	// Get the controller OwnerReference.
	owner := metav1.GetControllerOf(&existing)

	// Check if owner is nil.
	if owner == nil {
		return ctrl.Result{},
			fmt.Errorf("OwnerReference not found for SimpleApp %s/%s",
				app.Namespace,
				app.Name,
			)
	}

	// Check APIVersion, Kind and UID.
	if owner.APIVersion != cloudv1alpha1.GroupVersion.String() ||
		owner.Kind != "SimpleApp" ||
		owner.UID != app.UID {

		return ctrl.Result{},
			fmt.Errorf(
				"deployment %s/%s is not controlled by the current SimpleApp",
				app.Namespace,
				app.Name,
			)
	}

	// Only after ownership verification,
	// we can compare and update the Deployment.
	// Desired replicas and image
	desiredReplicas := app.Spec.Replicas
	desiredImage := app.Spec.Image
	changed := false

	// Compare Image and Replicas
	// If different, modify the certain target
	replicas := existing.Spec.Replicas

	if replicas == nil || *replicas != desiredReplicas {
		existing.Spec.Replicas = &desiredReplicas
		changed = true
	}

	containers := existing.Spec.Template.Spec.Containers

	// compare Image
	// to be sure the slice of containers is not 0
	if len(containers) == 0 {
		return ctrl.Result{}, fmt.Errorf(
			"Deployment %s/%s has no containers",
			existing.Namespace,
			existing.Name,
		)
	}

	if containers[0].Image != desiredImage {
		containers[0].Image = desiredImage
		changed = true
	}

	// If changed
	if changed == true {
		// execute Update
		if err := r.Client.Update(ctx, &existing); err != nil {
			return ctrl.Result{}, err
		}

		// Only finished Update the logger will be recorded
		logger.Info(
			"Deployment updated",
			"namespace", existing.Namespace,
			"name", existing.Name,
			"replicas", desiredReplicas,
			"image", desiredImage,
		)
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
