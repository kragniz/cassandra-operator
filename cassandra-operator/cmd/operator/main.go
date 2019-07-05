package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
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
	metricPollInterval   time.Duration
	metricRequestTimeout time.Duration
	logLevel             string
	allowEmptyDir        bool

	rootCmd = &cobra.Command{
		Use:               "cassandra-operator",
		Short:             "Operator for provisioning Cassandra clusters.",
		PersistentPreRunE: handleArgs,
		RunE:              startOperator,
	}

	scheme = runtime.NewScheme()
	log    = logf.Log.WithName("cassandra-operator")
)

func init() {
	v1alpha1.AddToScheme(scheme)
	kscheme.AddToScheme(scheme)

	rootCmd.PersistentFlags().DurationVar(&metricPollInterval, "metric-poll-interval", 5*time.Second, "Poll interval between cassandra nodes metrics retrieval")
	rootCmd.PersistentFlags().DurationVar(&metricRequestTimeout, "metric-request-timeout", 2*time.Second, "Time limit for cassandra node metrics requests")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "should be one of: debug, info, warn, error, fatal, panic")
	rootCmd.PersistentFlags().BoolVar(&allowEmptyDir, "allow-empty-dir", false, "Set to true in order to allow creation of clusters which use emptyDir storage")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func handleArgs(_ *cobra.Command, _ []string) error {
	var isPositive = func(duration time.Duration) bool {
		currentTime := time.Now()
		return currentTime.Add(duration).After(currentTime)
	}

	if !isPositive(metricPollInterval) {
		return fmt.Errorf("invalid metric-poll-interval, it must be a positive integer")
	}

	return nil
}

func startOperator(_ *cobra.Command, _ []string) error {
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
			client:             mgr.GetClient(),
			log:                log.WithName("reconciler"),
			receiver:           receiver,
			previousCassandras: map[string]*v1alpha1.Cassandra{},
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

	return nil
}
