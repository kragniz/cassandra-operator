module github.com/sky-uk/cassandra-operator/fake-cassandra-docker

go 1.12

require (
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	github.com/sky-uk/cassandra-operator/cassandra-operator v0.0.0-20190605131058-81ec3e5fbf77
	k8s.io/apimachinery v0.0.0-20190502092502-a44ef629a3c9
)

replace github.com/sky-uk/cassandra-operator/cassandra-operator => ../cassandra-operator

replace k8s.io/client-go => k8s.io/client-go v11.0.0+incompatible

replace k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190216013122-f05b8decd79c
