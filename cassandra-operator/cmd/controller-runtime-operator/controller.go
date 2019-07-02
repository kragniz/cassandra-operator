package main

import (
	"fmt"

	"github.com/go-logr/logr"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type reconcileCassandra struct {
	// client can be used to retrieve objects from the APIServer.
	client client.Client
	log    logr.Logger
}

// Implement reconcile.Reconciler so the controller can reconcile objects
var _ reconcile.Reconciler = &reconcileCassandra{}

func (r *reconcileCassandra) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	fmt.Println(request)
	return reconcile.Result{}, nil
}
