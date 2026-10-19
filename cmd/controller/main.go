package main

import (
	"log"

	"github.com/Leo-WANG17/pod-controller/internal/controller"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	cloudv1alpha1 "github.com/Leo-WANG17/pod-controller/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func main() {
	// get k8s cluster config, if no, then exit
	config, err := ctrl.GetConfig()
	if err != nil {
		log.Fatal("failed to get Kubernetes config: ", err)
	}

	// manager initialization
	// Create Scheme
	scheme := runtime.NewScheme()

	// Register K8s API types
	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		log.Fatal("failed to register Kubernetes API types: ", err)
	}

	// Register SimpleApp types
	if err := cloudv1alpha1.AddToScheme(scheme); err != nil {
		log.Fatal("failed to register SimpleApp API types: ", err)
	}

	mgr, err := ctrl.NewManager(
		config,
		ctrl.Options{
			Scheme: scheme,
		},
	)

	if err != nil {
		log.Fatalf("failed to create controller manager: %v", err)
	}

	// Client initialized by Manager
	k8sClient := mgr.GetClient()
	// Dependency injection
	podReconciler := &controller.PodReconciler{
		Client: k8sClient,
	}

	simpleAppReconciler := &controller.SimpleAppReconciler{
		Client: k8sClient,
	}
	// Register PodReconciler with Manager
	if err := podReconciler.SetupWithManager(mgr); err != nil {
		log.Fatal("reconciler initialization wrong", err)
		return
	}

	// Register SimpleAppReconciler with Manager
	if err := simpleAppReconciler.SetupWithManager(mgr); err != nil {
		log.Fatal("reconciler setup wrong", err)
		return
	}

	// Start the controller manager
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		log.Fatal("manager start failed", err)
	}
}
