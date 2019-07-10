package main

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/sky-uk/cassandra-operator/cassandra-operator/pkg/apis/cassandra/v1alpha1"
	v1alpha1helpers "github.com/sky-uk/cassandra-operator/cassandra-operator/pkg/apis/cassandra/v1alpha1/helpers"
	"github.com/sky-uk/cassandra-operator/cassandra-operator/pkg/apis/cassandra/v1alpha1/validation"
	"github.com/sky-uk/cassandra-operator/cassandra-operator/pkg/dispatcher"
	"github.com/sky-uk/cassandra-operator/cassandra-operator/pkg/operator/operations"
)

type reconcileCassandra struct {
	// client can be used to retrieve objects from the APIServer.
	client client.Client
	log    logr.Logger

	receiver           *operations.Receiver
	previousCassandras map[string]*v1alpha1.Cassandra
}

// Implement reconcile.Reconciler so the controller can reconcile objects
var _ reconcile.Reconciler = &reconcileCassandra{}

func (r *reconcileCassandra) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	fmt.Println(request)

	// set up a convinient log object so we don't have to type request over and over again
	log := r.log.WithValues("request", request)

	clusterID := fmt.Sprintf("%s.%s", request.Namespace, request.Name)

	// Fetch the Cassandra from the cache
	cass := &v1alpha1.Cassandra{}
	err := r.client.Get(context.TODO(), request.NamespacedName, cass)
	if errors.IsNotFound(err) {
		log.Error(nil, "Could not find Cassandra")

		deletedCass := &v1alpha1.Cassandra{ObjectMeta: metav1.ObjectMeta{Name: request.Name, Namespace: request.Namespace}}
		r.receiver.Receive(&dispatcher.Event{Kind: operations.DeleteCluster, Key: clusterID, Data: deletedCass})
		delete(r.previousCassandras, request.NamespacedName.String())
		return reconcile.Result{}, nil
	}

	if err != nil {
		log.Error(err, "Could not fetch Cassandra")
		return reconcile.Result{}, err
	}

	// Print the Cassandra
	log.Info("Reconciling Cassandra", "image name", cass.Spec.Pod.Image)

	v1alpha1helpers.SetDefaultsForCassandra(cass)
	err = validation.ValidateCassandra(cass).ToAggregate()
	if err != nil {
		log.Error(err, "validation error")
		return reconcile.Result{}, err
	}

	if cass.Annotations == nil {
		cass.Annotations = map[string]string{}
	}
	_, ok := cass.Annotations["reconciled.cassandra.core.sky.uk"]
	if !ok {
		// cassandra has not been created
		r.receiver.Receive(&dispatcher.Event{Kind: operations.AddCluster, Key: clusterID, Data: cass})
		cass.Annotations["reconciled.cassandra.core.sky.uk"] = "true"
	} else {
		previousCassandra, ok := r.previousCassandras[request.NamespacedName.String()]
		if !ok {
			return reconcile.Result{}, fmt.Errorf("couldn't find a previousCassandra")
		}
		r.receiver.Receive(&dispatcher.Event{
			Kind: operations.UpdateCluster,
			Key:  clusterID,
			Data: operations.ClusterUpdate{OldCluster: previousCassandra, NewCluster: cass},
		})
	}

	err = r.client.Update(context.TODO(), cass)
	if err != nil {
		log.Error(err, "Could not write Cassandra")
		return reconcile.Result{}, err
	}

	r.previousCassandras[request.NamespacedName.String()] = cass.DeepCopy()

	return reconcile.Result{}, nil
}
