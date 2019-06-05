package e2e_test

import (
	"os"
	"os/exec"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/sky-uk/cassandra-operator/cassandra-operator/pkg/apis/cassandra/v1alpha1"
	"github.com/sky-uk/cassandra-operator/cassandra-operator/pkg/cluster"
	"github.com/sky-uk/cassandra-operator/cassandra-operator/pkg/metrics"
	metricstesting "github.com/sky-uk/cassandra-operator/cassandra-operator/pkg/metrics/testing"
	"github.com/sky-uk/cassandra-operator/fake-cassandra-docker/test"
)

var (
	fakeCassandraImage string
	session            *Session
)

func init() {
	fakeCassandraImage = os.Getenv("FAKE_CASSANDRA_IMAGE")
	if fakeCassandraImage == "" {
		panic("FAKE_CASSANDRA_IMAGE must be supplied")
	}
}

func TestAll(t *testing.T) {
	BeforeSuite(func() {
		var err error
		command := exec.Command("docker", "run", "--rm", "--publish=7777:7777", fakeCassandraImage)
		session, err = Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ShouldNot(HaveOccurred())
		Eventually(session.Out, "5s").Should(Say("Starting fake Jolokia server"))
	})

	AfterSuite(func() {
		session.Interrupt()
		Eventually(session, "5s").Should(Exit())
	})

	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "e2e", test.CreateSequentialReporters("e2e"))
}

var _ = Context("Cassandra Docker container", func() {
	It("should accept Jolokia requests", func() {
		gathererer := metrics.NewGatherer(
			&metricstesting.StubbedJolokiaURLProvider{BaseURL: "http://localhost:7777"},
			&metrics.Config{RequestTimeout: 1 * time.Second},
		)
		clusterStatus, err := gathererer.GatherMetricsFor(
			aCluster("cluster1", "ns1"),
		)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(clusterStatus).ShouldNot(BeNil())
	})
})

func aCluster(clusterName, namespace string) *cluster.Cluster {
	theCluster, err := cluster.New(
		&v1alpha1.Cassandra{
			ObjectMeta: metav1.ObjectMeta{Name: clusterName, Namespace: namespace},
			Spec: v1alpha1.CassandraSpec{
				Racks: []v1alpha1.Rack{{Name: "a", Replicas: 1, StorageClass: "some-storage", Zone: "some-zone"}},
				Pod: v1alpha1.Pod{
					Memory:      resource.MustParse("1Gi"),
					CPU:         resource.MustParse("100m"),
					StorageSize: resource.MustParse("1Gi"),
				},
			},
		},
	)
	Expect(err).ToNot(HaveOccurred())
	return theCluster
}
