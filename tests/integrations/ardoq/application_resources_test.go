package ardoq_test

import (
	"K8SArdoqBridge/app/controllers"
	"K8SArdoqBridge/app/tests/helper"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

var _ = Describe("Deployments", Ordered, func() {
	var deploy *controllers.Resource
	var resourceName = helper.RandomString(10)
	var namespace = helper.RandomString(10)

	BeforeEach(func() {
		deploy = &controllers.Resource{
			Name:              resourceName,
			ResourceType:      "Deployment",
			Namespace:         namespace,
			Replicas:          helper.RandomInt(1, 5),
			Image:             "nginx-slim:0.8",
			CreationTimestamp: v1.Now().Format(time.RFC3339),
			Stack:             "nginx",
			Team:              "DevOps",
			Project:           "Test",
		}
		Expect(deploy.IsApplicationResourceValid()).To(BeTrue())
	})
	Context("Structural Validation", func() {
		It("Can get Deployment fields", func() {
			Expect(deploy.Name).To(Equal(resourceName))
			Expect(deploy.ResourceType).To(Equal("Deployment"))
			Expect(deploy.Namespace).To(Equal(namespace))
			Expect(deploy.Image).To(Equal("nginx-slim:0.8"))
			Expect(deploy.Replicas).To(Not(BeZero()))
		})
		It("Deployment name not set", func() {
			deploy.Name = ""
			Expect(deploy.IsApplicationResourceValid()).To(BeFalse())
		})
		It("Deployment namespace not set", func() {
			deploy.Namespace = ""
			Expect(deploy.IsApplicationResourceValid()).To(BeFalse())
		})
		It("Deployment image not set", func() {
			deploy.Image = ""
			Expect(deploy.IsApplicationResourceValid()).To(BeFalse())
		})
		It("Invalid resource Type set", func() {
			deploy.ResourceType = "Deplment"
			Expect(deploy.IsApplicationResourceValid()).To(BeFalse())
		})
	})
	Context("Deployment Ardoq Bridge tests", Ordered, func() {
		BeforeAll(func() {
			controllers.GenericUpsert("Namespace", namespace)
		})
		AfterAll(func() {
			err := controllers.GenericDelete("Namespace", namespace)
			Expect(err).ShouldNot(HaveOccurred())
			if err != nil {
				return
			}
		})
		It("Can create Deployment", func() {
			deploy.ID = controllers.GenericUpsert("Deployment", *deploy)
			Expect(deploy.ID).ShouldNot(BeNil())
		})
		It("Can Update Deployment", func() {
			deploy.Replicas += 1
			Expect(controllers.GenericUpsert("Deployment", *deploy)).ShouldNot(BeNil())
		})
		It("Can Delete Deployment", func() {
			Expect(controllers.GenericDelete("Deployment", *deploy)).Should(BeNil())
		})
		It("Can't Delete Non Existent Deployment", func() {
			deploy.Name = helper.RandomString(10)
			Expect(controllers.GenericDelete("Deployment", *deploy)).ShouldNot(BeNil())
		})
	})
})
var _ = Describe("StatefulSets", Ordered, func() {
	var sts *controllers.Resource
	var resourceName = helper.RandomString(10)
	var namespace = helper.RandomString(10)

	BeforeEach(func() {
		sts = &controllers.Resource{
			Name:              resourceName,
			ResourceType:      "StatefulSet",
			Namespace:         namespace,
			Replicas:          helper.RandomInt(1, 5),
			Image:             "postgresql:14.1",
			CreationTimestamp: v1.Now().Format(time.RFC3339),
			Stack:             "nginx",
			Team:              "DevOps",
			Project:           "Test",
		}
		Expect(sts.IsApplicationResourceValid()).To(BeTrue())
	})

	Context("Structural Validation", func() {
		It("Can get StatefulSet fields", func() {
			Expect(sts.Name).To(Equal(resourceName))
			Expect(sts.ResourceType).To(Equal("StatefulSet"))
			Expect(sts.Namespace).To(Equal(namespace))
			Expect(sts.Image).To(Equal("postgresql:14.1"))
			Expect(sts.Replicas).To(Not(BeZero()))
		})
		It("StatefulSet name not set", func() {
			sts.Name = ""
			Expect(sts.IsApplicationResourceValid()).To(BeFalse())
		})
		It("StatefulSet namespace not set", func() {
			sts.Namespace = ""
			Expect(sts.IsApplicationResourceValid()).To(BeFalse())
		})
		It("StatefulSet image not set", func() {
			sts.Image = ""
			Expect(sts.IsApplicationResourceValid()).To(BeFalse())
		})
		It("Invalid resource Type set", func() {
			sts.ResourceType = "Sttflst"
			Expect(sts.IsApplicationResourceValid()).To(BeFalse())
		})
	})
	Context("StatefulSet Ardoq Bridge tests", Ordered, func() {
		BeforeAll(func() {
			controllers.GenericUpsert("Namespace", namespace)
		})
		AfterAll(func() {
			err := controllers.GenericDelete("Namespace", namespace)
			Expect(err).ShouldNot(HaveOccurred())
			if err != nil {
				return
			}
		})
		It("Can create StatefulSet", func() {
			sts.ID = controllers.GenericUpsert("StatefulSet", *sts)
			Expect(sts.ID).ShouldNot(BeNil())
		})
		It("Can Update StatefulSet", func() {
			sts.Replicas += 1
			Expect(controllers.GenericUpsert("StatefulSet", *sts)).ShouldNot(BeNil())
		})
		It("Can Delete StatefulSet", func() {
			Expect(controllers.GenericDelete("StatefulSet", *sts)).Should(BeNil())
		})
		It("Can't Delete Non Existent StatefulSet", func() {
			sts.Name = helper.RandomString(10)
			Expect(controllers.GenericDelete("StatefulSet", *sts)).ShouldNot(BeNil())
		})
	})
})
