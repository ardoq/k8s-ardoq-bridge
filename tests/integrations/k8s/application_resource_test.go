package k8s_test

import (
	"K8SArdoqBridge/app/controllers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ApplicationResource", func() {
	Context("Deployment tests", func() {
		It("Can fetch created Deployments", func() {
			Expect(controllers.ApplicationResourceSearch("default", "Deployment", "nginx-deployment")).ShouldNot(BeNil())
		})
	})
	Context("StatefulSet tests", func() {
		It("Can fetch created StatefulSet", func() {
			Expect(controllers.ApplicationResourceSearch("default", "StatefulSet", "web")).ShouldNot(BeNil())
		})
	})
})

//var _ = Describe("Clusters", Ordered, func() {
//	var resourceName = helper.RandomString(10)
//	Context("Cluster Ardoq Link tests", Ordered, func() {
//		It("Can create Cluster", func() {
//			Expect(controllers.GenericUpsert("Cluster", resourceName)).ShouldNot(BeNil())
//			helper.ApplyDelay()
//		})
//		It("Can Lookup Cluster", func() {
//			Expect(controllers.LookupCluster(resourceName)).ShouldNot(BeNil())
//		})
//		It("Can Delete Cluster", func() {
//			Expect(controllers.GenericDelete("Cluster", resourceName)).Should(BeNil())
//		})
//
//		It("Can't Delete None Existent cluster", func() {
//			Expect(controllers.GenericDelete("Cluster", helper.RandomString(10))).ShouldNot(BeNil())
//		})
//	})
//})

//
//var _ = Describe("Deployments Manifests", func() {
//	//Context("Deployment tests", func() {
//	//	It("Can fetch created Deployments", func() {
//	//		helper.ApplyDelay()
//	//		Expect(controllers.ApplicationResourceSearch("default", "Deployment", "nginx-deployment")).ShouldNot(BeNil())
//	//	})
//	//})
//})

//var _ = Describe("Clusters", Ordered, func() {
//	var deploy *appsv1.Deployment
//	var resourceName = helper.RandomString(10)
//	var namespace = "default"
//	var ctx context.Context
//
//	BeforeEach(func() {
//		deploy = &appsv1.Deployment{
//			ObjectMeta: metav1.ObjectMeta{
//				Name:      resourceName,
//				Namespace: namespace,
//				Labels: map[string]string{
//					"sync-to-ardoq": "true",
//				},
//			},
//			Spec: appsv1.DeploymentSpec{
//				Replicas: pointer.Int32Ptr(2),
//				Selector: &metav1.LabelSelector{
//					MatchLabels: map[string]string{
//						"app": "test",
//					},
//				},
//				Template: v1.PodTemplateSpec{
//					ObjectMeta: metav1.ObjectMeta{
//						Labels: map[string]string{
//							"app": "test",
//						},
//					},
//					Spec: v1.PodSpec{
//						Containers: []v1.Container{
//							{
//								Name:  "web",
//								Image: "nginx:1.14.2",
//								Ports: []v1.ContainerPort{
//									{
//										Name:          "http",
//										Protocol:      v1.ProtocolTCP,
//										ContainerPort: 80,
//									},
//								},
//							},
//						},
//					},
//				},
//			},
//		}
//		err := inputs.Create(ctx, deploy)
//		Expect(err).ToNot(HaveOccurred())
//	})
//	Context("Deployment tests", Ordered, func() {
//		It("Can create Deployment", func() {
//			Expect(controllers.GenericUpsert("Deployment", *deploy)).ShouldNot(BeNil())
//			helper.ApplyDelay()
//		})
//		It("Can Update Deployment", func() {
//
//			Expect(controllers.GenericUpsert("Deployment", *deploy)).ShouldNot(BeNil())
//		})
//		It("Can Delete Deployment", func() {
//			Expect(controllers.GenericDelete("Deployment", *deploy)).Should(BeNil())
//		})
//	})
//})
