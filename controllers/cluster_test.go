package controllers_test

import (
	"ArdoqK8sBridge/app/controllers"
	"ArdoqK8sBridge/app/tests/helper"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Clusters", Ordered, func() {
	var resourceName = helper.RandomString(10)
	BeforeEach(func() {
	})
	Context("Cluster Ardoq Link tests", Ordered, func() {
		It("Can create Cluster", func() {
			Expect(controllers.UpsertCluster(resourceName)).ShouldNot(BeNil())
			helper.ApplyDelay()
		})
		It("Can Lookup Cluster", func() {
			Expect(controllers.LookupCluster(resourceName)).ShouldNot(BeNil())
		})
		It("Can Delete Cluster", func() {
			Expect(controllers.DeleteCluster(resourceName)).Should(BeNil())
		})

		It("Can't Delete None Existent cluster", func() {
			Expect(controllers.DeleteCluster(helper.RandomString(10))).ShouldNot(BeNil())
		})
	})
})
