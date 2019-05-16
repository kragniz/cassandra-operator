package helpers

import (
	"github.com/sky-uk/cassandra-operator/cassandra-operator/pkg/apis/cassandra/v1alpha1"
)

// UseEmptyDir returns a dereferenced value for Spec.UseEmptyDir
func UseEmptyDir(c *v1alpha1.Cassandra) bool {
	if c.Spec.UseEmptyDir != nil {
		return *c.Spec.UseEmptyDir
	}
	return false
}
