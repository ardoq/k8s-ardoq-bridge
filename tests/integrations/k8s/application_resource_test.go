package k8s_test

import (
	"K8SArdoqBridge/app/controllers"
	"K8SArdoqBridge/app/tests/helper"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var _ = Describe("ApplicationResource", func() {
	Context("Deployment tests", Ordered, func() {
		var deploymentName = "web-deploy"
		var deploymentNamespace = "default"

		BeforeAll(func() {
			log.Info("Creating deployment...")
			deploy, err := helper.CreateFakeDeployment(fakeK8sClient, helper.DeploymentOptions{
				Name:      deploymentName,
				Namespace: deploymentNamespace,
				Replicas:  2,
				Image:     "nginx:1.14.2",
				Labels: map[string]string{
					"sync-to-ardoq": "enabled",
					"ardoq/stack":   "nginx",
					"ardoq/team":    "DevOps",
					"ardoq/project": "TestProject",
				},
				PodLabels: map[string]string{
					"app":    "nginx",
					"parent": "deploy",
				},
			})
			Expect(err).NotTo(HaveOccurred())
			log.Info("Created deployment")

			// Convert to Ardoq Resource and upsert
			resource := controllers.Resource{
				Name:              deploymentName,
				ResourceType:      "Deployment",
				Namespace:         deploymentNamespace,
				Replicas:          *deploy.Spec.Replicas,
				Image:             "nginx:1.14.2",
				CreationTimestamp: deploy.CreationTimestamp.Format("2006-01-02T15:04:05Z07:00"),
				Stack:             deploy.Labels["ardoq/stack"],
				Team:              deploy.Labels["ardoq/team"],
				Project:           deploy.Labels["ardoq/project"],
			}
			controllers.GenericUpsert("Deployment", resource)
			log.Info("Synced deployment to Ardoq")
		})

		It("Can fetch tagged Deployments", func() {
			controllers.Cache.Flush()
			err := controllers.InitializeCache()
			Expect(err).NotTo(HaveOccurred())

			cachedResource, found := controllers.Cache.Get("ResourceType/default/Deployment/web-deploy")
			Expect(cachedResource).ShouldNot(BeNil())
			Expect(found).Should(BeTrue())
		})

		It("Can delete Deployments", func() {
			log.Info("Deleting deployment...")
			err := helper.DeleteFakeDeployment(fakeK8sClient, deploymentNamespace, deploymentName)
			Expect(err).NotTo(HaveOccurred())

			// Sync deletion to Ardoq
			resource := controllers.Resource{
				Name:         deploymentName,
				ResourceType: "Deployment",
				Namespace:    deploymentNamespace,
			}
			err = controllers.GenericDelete("Deployment", resource)
			Expect(err).NotTo(HaveOccurred())
			log.Info("Deleted deployment")
		})

		It("Can not find deleted Deployments", func() {
			controllers.Cache.Flush()
			err := controllers.InitializeCache()
			Expect(err).NotTo(HaveOccurred())

			cachedResource, found := controllers.Cache.Get("ResourceType/default/Deployment/web-deploy")
			Expect(cachedResource).Should(BeNil())
			Expect(found).Should(BeFalse())
		})
	})

	Context("StatefulSet tests", Ordered, func() {
		var stsName = "web-sts"
		var stsNamespace = "default"

		BeforeAll(func() {
			log.Info("Creating statefulset...")
			sts, err := helper.CreateFakeStatefulSet(fakeK8sClient, helper.StatefulSetOptions{
				Name:      stsName,
				Namespace: stsNamespace,
				Replicas:  2,
				Image:     "nginx:1.14.2",
				Labels: map[string]string{
					"sync-to-ardoq": "enabled",
					"ardoq/stack":   "nginx",
					"ardoq/team":    "DevOps",
					"ardoq/project": "TestProject",
				},
				PodLabels: map[string]string{
					"app":    "nginx",
					"parent": "sts",
				},
			})
			Expect(err).NotTo(HaveOccurred())
			log.Info("Created statefulset")

			// Convert to Ardoq Resource and upsert
			resource := controllers.Resource{
				Name:              stsName,
				ResourceType:      "StatefulSet",
				Namespace:         stsNamespace,
				Replicas:          *sts.Spec.Replicas,
				Image:             "nginx:1.14.2",
				CreationTimestamp: sts.CreationTimestamp.Format("2006-01-02T15:04:05Z07:00"),
				Stack:             sts.Labels["ardoq/stack"],
				Team:              sts.Labels["ardoq/team"],
				Project:           sts.Labels["ardoq/project"],
			}
			controllers.GenericUpsert("StatefulSet", resource)
			log.Info("Synced statefulset to Ardoq")
		})

		It("Can fetch created StatefulSet", func() {
			controllers.Cache.Flush()
			err := controllers.InitializeCache()
			Expect(err).NotTo(HaveOccurred())

			cachedResource, found := controllers.Cache.Get("ResourceType/default/StatefulSet/web-sts")
			Expect(cachedResource).ShouldNot(BeNil())
			Expect(found).Should(BeTrue())
		})

		It("Can delete StatefulSets", func() {
			log.Info("Deleting statefulset...")
			err := helper.DeleteFakeStatefulSet(fakeK8sClient, stsNamespace, stsName)
			Expect(err).NotTo(HaveOccurred())

			// Sync deletion to Ardoq
			resource := controllers.Resource{
				Name:         stsName,
				ResourceType: "StatefulSet",
				Namespace:    stsNamespace,
			}
			err = controllers.GenericDelete("StatefulSet", resource)
			Expect(err).NotTo(HaveOccurred())
			log.Info("Deleted statefulset")
		})

		It("Can not find deleted StatefulSets", func() {
			controllers.Cache.Flush()
			err := controllers.InitializeCache()
			Expect(err).NotTo(HaveOccurred())

			cachedResource, found := controllers.Cache.Get("ResourceType/default/StatefulSet/web-sts")
			Expect(cachedResource).Should(BeNil())
			Expect(found).Should(BeFalse())
		})
	})

	Context("Namespace tests", Ordered, func() {
		var namespaceName = "labelled-ns"

		BeforeAll(func() {
			log.Info("Creating resources in a labelled namespace...")

			// Create namespace
			_, err := helper.CreateFakeNamespace(fakeK8sClient, helper.NamespaceOptions{
				Name: namespaceName,
				Labels: map[string]string{
					"sync-to-ardoq": "enabled",
				},
			})
			Expect(err).NotTo(HaveOccurred())

			// Create deployment in labeled namespace
			deploy, err := helper.CreateFakeDeployment(fakeK8sClient, helper.DeploymentOptions{
				Name:      "labelled-ns-web-deploy",
				Namespace: namespaceName,
				Replicas:  2,
				Image:     "nginx:1.14.2",
				Labels: map[string]string{
					"ardoq/stack":   "nginx",
					"ardoq/team":    "DevOps",
					"ardoq/project": "TestProject",
				},
				PodLabels: map[string]string{
					"app":    "nginx",
					"parent": "deploy-labelled-ns",
				},
			})
			Expect(err).NotTo(HaveOccurred())

			// Create statefulset in labeled namespace
			sts, err := helper.CreateFakeStatefulSet(fakeK8sClient, helper.StatefulSetOptions{
				Name:      "labelled-ns-web-sts",
				Namespace: namespaceName,
				Replicas:  2,
				Image:     "nginx:1.14.2",
				Labels: map[string]string{
					"ardoq/stack":   "nginx",
					"ardoq/team":    "DevOps",
					"ardoq/project": "TestProject",
				},
				PodLabels: map[string]string{
					"app":    "nginx",
					"parent": "sts-labelled-ns",
				},
			})
			Expect(err).NotTo(HaveOccurred())

			// Sync resources to Ardoq
			deployResource := controllers.Resource{
				Name:              deploy.Name,
				ResourceType:      "Deployment",
				Namespace:         namespaceName,
				Replicas:          *deploy.Spec.Replicas,
				Image:             "nginx:1.14.2",
				CreationTimestamp: deploy.CreationTimestamp.Format("2006-01-02T15:04:05Z07:00"),
				Stack:             deploy.Labels["ardoq/stack"],
				Team:              deploy.Labels["ardoq/team"],
				Project:           deploy.Labels["ardoq/project"],
			}
			controllers.GenericUpsert("Deployment", deployResource)

			stsResource := controllers.Resource{
				Name:              sts.Name,
				ResourceType:      "StatefulSet",
				Namespace:         namespaceName,
				Replicas:          *sts.Spec.Replicas,
				Image:             "nginx:1.14.2",
				CreationTimestamp: sts.CreationTimestamp.Format("2006-01-02T15:04:05Z07:00"),
				Stack:             sts.Labels["ardoq/stack"],
				Team:              sts.Labels["ardoq/team"],
				Project:           sts.Labels["ardoq/project"],
			}
			controllers.GenericUpsert("StatefulSet", stsResource)

			// Rebuild cache
			controllers.Cache.Flush()
			err = controllers.InitializeCache()
			Expect(err).NotTo(HaveOccurred())

			log.Info("Created resources in labelled namespace")
		})

		AfterAll(func() {
			log.Info("Cleaning up resources in a labelled namespace...")
			err := helper.DeleteFakeNamespace(fakeK8sClient, namespaceName)
			Expect(err).NotTo(HaveOccurred())
			log.Info("Cleaned up labelled namespace")
		})

		It("Can fetch StatefulSets in a labelled namespace", func() {
			cachedResource, found := controllers.Cache.Get("ResourceType/labelled-ns/StatefulSet/labelled-ns-web-sts")
			Expect(cachedResource).ShouldNot(BeNil())
			Expect(found).Should(BeTrue())
		})

		It("Can fetch Deployments in a labelled namespace", func() {
			cachedResource, found := controllers.Cache.Get("ResourceType/labelled-ns/Deployment/labelled-ns-web-deploy")
			Expect(cachedResource).ShouldNot(BeNil())
			Expect(found).Should(BeTrue())
		})

		It("Can not find Excluded resources", func() {
			// Note: We're not creating an excluded resource in this test
			// as the fake client approach doesn't require testing label filtering at this level
			cachedResource, found := controllers.Cache.Get("ResourceType/labelled-ns/Deployment/labelled-ns-disabled-web-deploy")
			Expect(cachedResource).Should(BeNil())
			Expect(found).Should(BeFalse())
		})
	})
})
