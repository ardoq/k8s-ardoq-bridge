package ardoq_test

import (
	"K8SArdoqBridge/app/controllers"
	"K8SArdoqBridge/app/tests/helper"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Clusters", Ordered, func() {
	var resourceName = helper.RandomString(10)
	BeforeEach(func() {
	})
	Context("Cluster Ardoq Link tests", Ordered, func() {
		It("Can create Cluster", func() {
			Expect(controllers.GenericUpsert("Cluster", resourceName)).ShouldNot(BeNil())
			helper.ApplyDelay()
		})
		It("Can Lookup Cluster", func() {
			Expect(controllers.LookupCluster(resourceName)).ShouldNot(BeNil())
		})
		It("Can Delete Cluster", func() {
			Expect(controllers.GenericDelete("Cluster", resourceName)).Should(BeNil())
		})

		It("Can't Delete None Existent cluster", func() {
			Expect(controllers.GenericDelete("Cluster", helper.RandomString(10))).ShouldNot(BeNil())
		})
	})
})
