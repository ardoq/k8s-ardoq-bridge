package ardoq_test

import (
	"K8SArdoqBridge/app/controllers"
	"K8SArdoqBridge/app/tests/helper"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Namespaces", Ordered, func() {
	var resourceName = helper.RandomString(10)
	Context("Namespace Ardoq Link tests", Ordered, func() {
		It("Can create Namespace", func() {
			Expect(controllers.GenericUpsert("Namespace", resourceName)).ShouldNot(BeNil())
		})
		It("Can Delete Namespace", func() {
			Expect(controllers.GenericDelete("Namespace", resourceName)).Should(Succeed())
		})
		It("Can't Delete Non Existent Namespace", func() {
			Expect(controllers.GenericDelete("Namespace", helper.RandomString(10))).ShouldNot(Succeed())
		})
	})
})
