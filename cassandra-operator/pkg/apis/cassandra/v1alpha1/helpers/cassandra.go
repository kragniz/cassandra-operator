package helpers

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/sky-uk/cassandra-operator/cassandra-operator/pkg/apis/cassandra"
	"github.com/sky-uk/cassandra-operator/cassandra-operator/pkg/apis/cassandra/v1alpha1"
)

func NewControllerRef(c *v1alpha1.Cassandra) metav1.OwnerReference {
	return *metav1.NewControllerRef(c, schema.GroupVersionKind{
		Group:   cassandra.GroupName,
		Version: cassandra.Version,
		Kind:    cassandra.Kind,
	})
}

// UseEmptyDir returns a dereferenced value for Spec.UseEmptyDir
func UseEmptyDir(c *v1alpha1.Cassandra) bool {
	if c.Spec.UseEmptyDir != nil {
		return *c.Spec.UseEmptyDir
	}
	return false
}

// GetImage returns the image for a cluster
func GetCassandraImage(c *v1alpha1.Cassandra) string {
	if c.Spec.Pod.Image != nil {
		return *c.Spec.Pod.Image
	}
	return v1alpha1.DefaultCassandraImage
}

// GetBootstrapperImage returns the bootstrapper image for a cluster
func GetBootstrapperImage(c *v1alpha1.Cassandra) string {
	if c.Spec.Pod.BootstrapperImage != nil {
		return *c.Spec.Pod.BootstrapperImage
	}
	return v1alpha1.DefaultCassandraBootstrapperImage
}

// GetSnapshopImage returns the snapshot image for a cluster
func GetSnapshopImage(c *v1alpha1.Cassandra) string {
	if c.Spec.Snapshot != nil {
		if c.Spec.Snapshot.Image != nil {
			return *c.Spec.Snapshot.Image
		}
	}
	return v1alpha1.DefaultCassandraSnapshotImage
}

func GetDatacenter(c *v1alpha1.Cassandra) string {
	if c.Spec.Datacenter == nil {
		return v1alpha1.DefaultDCName
	}
	return *c.Spec.Datacenter
}
