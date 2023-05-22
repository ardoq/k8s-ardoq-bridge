package k8s_test

import (
	"K8SArdoqBridge/app/controllers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

var _ = Describe("ApplicationResource", func() {
	Context("Deployment tests", Ordered, func() {
		BeforeAll(func() {
			log.Info("Creating deployment...")
			cmd := exec.Command("kubectl", "apply", "--wait=true", "-f", "manifests/deployment.yaml")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session.Out, 5).Should(gbytes.Say(".*deployment.apps.* [created|unchanged|configured].*"))
			cmd = exec.Command("kubectl", "wait", "--for=condition=ready", "--timeout=180s", "pod", "-l", "app=nginx,parent=deploy")
			session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session.Out, 20).Should(gbytes.Say(".*pod.* met*"))
			log.Infof("Created deployment")
		})
		It("Can fetch tagged Deployments", func() {
			controllers.Cache.Flush()
			err := controllers.InitializeCache()
			if err != nil {
				log.Fatalf("Error rebuilding cache: %s", err.Error())
			}
			cachedResource, found := controllers.Cache.Get("ResourceType/default/Deployment/web-deploy")
			Expect(cachedResource).ShouldNot(BeNil())
			Expect(found).Should(BeTrue())
		})
		It("Can delete Deployments", func() {
			log.Info("Deleting deployment...")
			cmd := exec.Command("kubectl", "delete", "--wait=true", "-f", "manifests/deployment.yaml")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session.Out, 5).Should(gbytes.Say(".*deployment.apps.* deleted.*"))

			cmd = exec.Command("kubectl", "wait", "--for=delete", "--timeout=180s", "pod", "-l", "app=nginx,parent=deploy")
			session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			log.Info("Deleted deployment.")
		})
		It("Can not find deleted Deployments", func() {
			Eventually(session.Err, 20).Should(gbytes.Say(`.*Deleted Deployment: web-deploy*`))

			controllers.Cache.Flush()
			err := controllers.InitializeCache()
			if err != nil {
				log.Fatalf("Error rebuilding cache: %s", err.Error())
			}
			cachedResource, found := controllers.Cache.Get("ResourceType/default/Deployment/web-deploy")
			Expect(cachedResource).Should(BeNil())
			Expect(found).Should(BeFalse())
		})
	})

	Context("StatefulSet tests", Ordered, func() {
		BeforeAll(func() {
			log.Info("Creating statefulset...")
			cmd := exec.Command("kubectl", "apply", "--wait=true", "-f", "manifests/statefulset.yaml")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session.Out, 5).Should(gbytes.Say(".*statefulset.apps.* [created|unchanged|configured].*"))
			cmd = exec.Command("kubectl", "wait", "--for=condition=ready", "--timeout=180s", "pod", "-l", "app=nginx,parent=sts")
			session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session.Out, 20).Should(gbytes.Say(".*pod.* met*"))
			log.Infof("Created statefulset")
		})
		It("Can fetch created StatefulSet", func() {
			controllers.Cache.Flush()
			err := controllers.InitializeCache()
			if err != nil {
				log.Fatalf("Error rebuilding cache: %s", err.Error())
			}
			cachedResource, found := controllers.Cache.Get("ResourceType/default/StatefulSet/web-sts")
			Expect(cachedResource).ShouldNot(BeNil())
			Expect(found).Should(BeTrue())
		})
		It("Can deleted StatefulSets", func() {
			log.Info("Deleting statefulset...")
			cmd := exec.Command("kubectl", "delete", "--wait=true", "-f", "manifests/statefulset.yaml")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session.Out, 5).Should(gbytes.Say(".*statefulset.apps.* deleted.*"))

			cmd = exec.Command("kubectl", "wait", "--for=delete", "--timeout=180s", "pod", "-l", "app=nginx,parent=sts")
			session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			log.Info("Deleted statefulset.")
		})
		It("Can not find deleted StatefulSets", func() {
			Eventually(session.Err, 20).Should(gbytes.Say(`.*Deleted StatefulSet: web-sts*`))

			controllers.Cache.Flush()
			err := controllers.InitializeCache()
			if err != nil {
				log.Fatalf("Error rebuilding cache: %s", err.Error())
			}
			cachedResource, found := controllers.Cache.Get("ResourceType/default/StatefulSet/web-sts")
			Expect(cachedResource).Should(BeNil())
			Expect(found).Should(BeFalse())
		})
	})
	Context("Namespace tests", Ordered, func() {
		BeforeAll(func() {
			log.Info("Creating resources in a labelled namespace...")
			cmd := exec.Command("kubectl", "apply", "--wait=true", "-Rf", "manifests/labeled-ns/")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session.Out, 5).Should(gbytes.Say(".*deployment.apps.* [created|unchanged|configured].*"))
			Eventually(session.Out, 5).Should(gbytes.Say(".*statefulset.apps.* [created|unchanged|configured].*"))

			cmd = exec.Command("kubectl", "wait", "--for=condition=ready", "--timeout=180s", "pod", "-l", "parent=sts-labelled-ns", "-n", "labelled-ns")
			stsSession, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(stsSession.Out, 20).Should(gbytes.Say(".*pod.* met*"))

			cmd = exec.Command("kubectl", "wait", "--for=condition=ready", "--timeout=180s", "pod", "-l", "parent=deploy-labelled-ns", "-n", "labelled-ns")
			deploySession, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(deploySession.Out, 10).Should(gbytes.Say(".*pod.* met*"))

			cmd = exec.Command("kubectl", "wait", "--for=condition=ready", "--timeout=180s", "pod", "-l", "parent=deploy-excluded", "-n", "labelled-ns")
			excludedSession, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(excludedSession.Out, 10).Should(gbytes.Say(".*pod.* met*"))

			controllers.Cache.Flush()
			err = controllers.InitializeCache()
			if err != nil {
				log.Fatalf("Error rebuilding cache: %s", err.Error())
			}
		})
		AfterAll(func() {
			log.Info("Cleaning up resources in a labelled namespace...")
			cmd := exec.Command("kubectl", "delete", "--wait=true", "-Rf", "manifests/labeled-ns/")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session.Out, 5).Should(gbytes.Say(".*deleted.*"))
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
			cachedResource, found := controllers.Cache.Get("ResourceType/labelled-ns/StatefulSet/labelled-ns-disbaled-web-deploy")
			Expect(cachedResource).Should(BeNil())
			Expect(found).Should(BeFalse())
		})
	})
})
