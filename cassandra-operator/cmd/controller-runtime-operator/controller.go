package main

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"

	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/sky-uk/cassandra-operator/cassandra-operator/pkg/apis/cassandra/v1alpha1"
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

	// set up a convinient log object so we don't have to type request over and over again
	log := r.log.WithValues("request", request)

	// Fetch the Cassandra from the cache
	cass := &v1alpha1.Cassandra{}
	err := r.client.Get(context.TODO(), request.NamespacedName, cass)
	if errors.IsNotFound(err) {
		log.Error(nil, "Could not find Cassandra")
		return reconcile.Result{}, nil
	}

	if err != nil {
		log.Error(err, "Could not fetch Cassandra")
		return reconcile.Result{}, err
	}

	// Print the Cassandra
	log.Info("Reconciling Cassandra", "image name", cass.Spec.Pod.Image)

	// Set the label if it is missing
	if cass.Labels == nil {
		cass.Labels = map[string]string{}
	}
	if cass.Labels["hello"] == "world" {
		return reconcile.Result{}, nil
	}

	// Update the Cassandra
	cass.Labels["hello"] = "world"
	err = r.client.Update(context.TODO(), cass)
	if err != nil {
		log.Error(err, "Could not write Cassandra")
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
