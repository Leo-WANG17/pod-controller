package main

import (
	"log"

	"github.com/Leo-WANG17/pod-controller/internal/controller"

	ctrl "sigs.k8s.io/controller-runtime"
)

func main() {
	// get k8s cluster config, if no, then exit
	config, err := ctrl.GetConfig()
	if err != nil {
		log.Fatal("failed to get Kubernetes config: ", err)
	}

	// manager initialization
	mgr, err := ctrl.NewManager(
		config,
		ctrl.Options{},
	)

	if err != nil {
		log.Fatal("manager initialization wrong", err)
		return
	}

	// Client initialized by Manager
	client := mgr.GetClient()
	// Dependency injection
	podReconciler := &controller.PodReconciler{
		Client: client,
	}
	// Register PodReconciler with Manager
	if err := podReconciler.SetupWithManager(mgr); err != nil {
		log.Fatal("reconciler initialization wrong", err)
		return
	}

	// Start the controller manager
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		log.Fatal("manager start failed", err)
		return
	}
}
