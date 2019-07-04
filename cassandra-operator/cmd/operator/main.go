package main

import (
	"flag"
	"os"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	kscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/sky-uk/cassandra-operator/cassandra-operator/pkg/apis/cassandra/v1alpha1"
	"github.com/sky-uk/cassandra-operator/cassandra-operator/pkg/client/clientset/versioned"
	"github.com/sky-uk/cassandra-operator/cassandra-operator/pkg/cluster"
	"github.com/sky-uk/cassandra-operator/cassandra-operator/pkg/operator/operations"
)

var (
	scheme = runtime.NewScheme()
	log    = logf.Log.WithName("cassandra-operator")
)

func init() {
	v1alpha1.AddToScheme(scheme)
	kscheme.AddToScheme(scheme)
}

func main() {
	flag.Parse()
	logf.SetLogger(zap.Logger(false))
	entryLog := log.WithName("entrypoint")

	kubeConfig := config.GetConfigOrDie()

	kubeClientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		entryLog.Error(err, "Unable to obtain clientset")
		os.Exit(1)
	}

	cassandraClientset := versioned.NewForConfigOrDie(kubeConfig)

	// Setup a Manager
	entryLog.Info("setting up manager")
	mgr, err := manager.New(kubeConfig, manager.Options{Scheme: scheme})
	if err != nil {
		entryLog.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}

	clusters := make(map[string]*cluster.Cluster)

	eventRecorder := cluster.NewEventRecorder(kubeClientset)
	clusterAccessor := cluster.NewAccessor(kubeClientset, cassandraClientset, eventRecorder)

	receiver := operations.NewEventReceiver(
		clusters,
		clusterAccessor,
		nil,
		eventRecorder,
	)

	// Setup a new controller to reconcile ReplicaSets
	entryLog.Info("Setting up controller")
	c, err := controller.New("cassandra", mgr, controller.Options{
		Reconciler: &reconcileCassandra{
			previousCassandras: map[string]*v1alpha1.Cassandra{},
			client:             mgr.GetClient(),
			log:                log.WithName("reconciler"),
			receiver:           receiver,
		},
	})
	if err != nil {
		entryLog.Error(err, "unable to set up individual controller")
		os.Exit(1)
	}

	// Watch Cassandras and enqueue Cassandra object key
	if err := c.Watch(&source.Kind{Type: &v1alpha1.Cassandra{}}, &handler.EnqueueRequestForObject{}); err != nil {
		entryLog.Error(err, "unable to watch Cassandras")
		os.Exit(1)
	}

	// Watch StatefulSets and enqueue owning Cassandra key
	if err := c.Watch(&source.Kind{Type: &appsv1.StatefulSet{}},
		&handler.EnqueueRequestForOwner{OwnerType: &v1alpha1.Cassandra{}, IsController: true}); err != nil {
		entryLog.Error(err, "unable to watch StatefulSets")
		os.Exit(1)
	}

	// Watch Pods and enqueue owning Cassandra key
	if err := c.Watch(&source.Kind{Type: &corev1.Pod{}},
		&handler.EnqueueRequestForOwner{OwnerType: &v1alpha1.Cassandra{}, IsController: true}); err != nil {
		entryLog.Error(err, "unable to watch Pods")
		os.Exit(1)
	}

	entryLog.Info("starting manager")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		entryLog.Error(err, "unable to run manager")
		os.Exit(1)
	}
}
