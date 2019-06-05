package e2e_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/sky-uk/cassandra-operator/fake-cassandra-docker/test"
)

func TestAll(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "My Suite")
	RunSpecsWithDefaultAndCustomReporters(t, "E2E Suite", test.CreateSequentialReporters("e2e"))
}

var _ = Context("some context", func() {
	It("it should foo", func() {
		By("doing something")
		Expect("foo").Should(Equal("bar"))
	})
})
