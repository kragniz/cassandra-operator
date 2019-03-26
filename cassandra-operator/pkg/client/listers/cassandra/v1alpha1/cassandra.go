// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/sky-uk/cassandra-operator/cassandra-operator/pkg/apis/cassandra/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// CassandraLister helps list Cassandras.
type CassandraLister interface {
	// List lists all Cassandras in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.Cassandra, err error)
	// Cassandras returns an object that can list and get Cassandras.
	Cassandras(namespace string) CassandraNamespaceLister
	CassandraListerExpansion
}

// cassandraLister implements the CassandraLister interface.
type cassandraLister struct {
	indexer cache.Indexer
}

// NewCassandraLister returns a new CassandraLister.
func NewCassandraLister(indexer cache.Indexer) CassandraLister {
	return &cassandraLister{indexer: indexer}
}

// List lists all Cassandras in the indexer.
func (s *cassandraLister) List(selector labels.Selector) (ret []*v1alpha1.Cassandra, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Cassandra))
	})
	return ret, err
}

// Cassandras returns an object that can list and get Cassandras.
func (s *cassandraLister) Cassandras(namespace string) CassandraNamespaceLister {
	return cassandraNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// CassandraNamespaceLister helps list and get Cassandras.
type CassandraNamespaceLister interface {
	// List lists all Cassandras in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.Cassandra, err error)
	// Get retrieves the Cassandra from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.Cassandra, error)
	CassandraNamespaceListerExpansion
}

// cassandraNamespaceLister implements the CassandraNamespaceLister
// interface.
type cassandraNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Cassandras in the indexer for a given namespace.
func (s cassandraNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.Cassandra, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Cassandra))
	})
	return ret, err
}

// Get retrieves the Cassandra from the indexer for a given namespace and name.
func (s cassandraNamespaceLister) Get(name string) (*v1alpha1.Cassandra, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("cassandra"), name)
	}
	return obj.(*v1alpha1.Cassandra), nil
}
