package controllers_test

import (
	"ArdoqK8sBridge/app/controllers"
	"ArdoqK8sBridge/app/tests/helper"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Namespaces", Ordered, func() {
	var resourceName = helper.RandomString(10)
	BeforeEach(func() {
	})
	Context("Namespace Ardoq Link tests", Ordered, func() {
		It("Can create Namespace", func() {
			Expect(controllers.UpsertNamespace(resourceName)).ShouldNot(BeNil())
			helper.ApplyDelay()
		})
		It("Can Delete Namespace", func() {
			Expect(controllers.DeleteNamespace(resourceName)).Should(BeNil())
		})
		It("Can't Delete None Existent Namespace", func() {
			Expect(controllers.DeleteNamespace(helper.RandomString(10))).ShouldNot(BeNil())
		})
	})
})
