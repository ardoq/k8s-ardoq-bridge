package k8s_test

import (
	"K8SArdoqBridge/app/controllers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"k8s.io/klog/v2"
	"os/exec"
)

var _ = Describe("ApplicationResource", func() {
	Context("Deployment tests", Ordered, func() {
		BeforeAll(func() {
			klog.Info("Creating deployment...")
			cmd := exec.Command("kubectl", "apply", "--wait=true", "-f", "manifests/deployment.yaml")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session.Out, 5).Should(gbytes.Say(".*deployment.apps.* [created|unchanged|configured].*"))
			cmd = exec.Command("kubectl", "wait", "--for=condition=ready", "--timeout=180s", "pod", "-l", "app=nginx,parent=deploy")
			session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session.Out, 10).Should(gbytes.Say(".*pod.* met*"))
			klog.Infof("Created deployment")
		})
		It("Can fetch tagged Deployments", func() {
			Eventually(session.Err, 20).Should(gbytes.Say(`.*Added Deployment: "web-deploy"*`))

			//re-initialize cache and confirm content
			controllers.Cache.Flush()
			err := controllers.InitializeCache()
			if err != nil {
				klog.Fatalf("Error rebuilding cache: %s", err.Error())
			}
			cachedResource, _ := controllers.Cache.Get("ResourceType/default/Deployment/web-deploy")
			Expect(cachedResource).ShouldNot(BeNil())
		})
		It("Can delete Deployments", func() {
			klog.Info("Deleting deployment...")
			cmd := exec.Command("kubectl", "delete", "--wait=true", "-f", "manifests/deployment.yaml")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session.Out, 5).Should(gbytes.Say(".*deployment.apps.* deleted.*"))

			cmd = exec.Command("kubectl", "wait", "--for=delete", "--timeout=180s", "pod", "-l", "app=nginx,parent=deploy")
			session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			klog.Info("Deleted deployment.")
		})
		It("Can not find deleted Deployments", func() {
			Eventually(session.Err, 20).Should(gbytes.Say(`.*Deleted Deployment: "web-deploy"*`))

			//re-initialize cache and confirm content
			controllers.Cache.Flush()
			err := controllers.InitializeCache()
			if err != nil {
				klog.Fatalf("Error rebuilding cache: %s", err.Error())
			}
			cachedResource, _ := controllers.Cache.Get("ResourceType/default/Deployment/web-deploy")
			Expect(cachedResource).Should(BeNil())
		})
	})

	Context("StatefulSet tests", Ordered, func() {
		BeforeAll(func() {
			klog.Info("Creating statefulset...")
			cmd := exec.Command("kubectl", "apply", "--wait=true", "-f", "manifests/statefulset.yaml")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session.Out, 5).Should(gbytes.Say(".*statefulset.apps.* [created|unchanged|configured].*"))
			cmd = exec.Command("kubectl", "wait", "--for=condition=ready", "--timeout=180s", "pod", "-l", "app=nginx,parent=sts")
			session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session.Out, 10).Should(gbytes.Say(".*pod.* met*"))
			klog.Infof("Created statefulset")
		})
		It("Can fetch created StatefulSet", func() {
			Eventually(session.Err, 20).Should(gbytes.Say(`.*Added StatefulSet: "web-sts"*`))

			//re-initialize cache and confirm content
			controllers.Cache.Flush()
			err := controllers.InitializeCache()
			if err != nil {
				klog.Fatalf("Error rebuilding cache: %s", err.Error())
			}
			cachedResource, _ := controllers.Cache.Get("ResourceType/default/StatefulSet/web-sts")
			Expect(cachedResource).ShouldNot(BeNil())
		})
		It("Can deleted StatefulSets", func() {
			klog.Info("Deleting statefulset...")
			cmd := exec.Command("kubectl", "delete", "--wait=true", "-f", "manifests/statefulset.yaml")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session.Out, 5).Should(gbytes.Say(".*statefulset.apps.* deleted.*"))

			cmd = exec.Command("kubectl", "wait", "--for=delete", "--timeout=180s", "pod", "-l", "app=nginx,parent=sts")
			session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			klog.Info("Deleted statefulset.")
		})
		It("Can not find deleted StatefulSets", func() {
			Eventually(session.Err, 20).Should(gbytes.Say(`.*Deleted StatefulSet: "web-sts"*`))

			//re-initialize cache and confirm content
			controllers.Cache.Flush()
			err := controllers.InitializeCache()
			if err != nil {
				klog.Fatalf("Error rebuilding cache: %s", err.Error())
			}
			cachedResource, _ := controllers.Cache.Get("ResourceType/default/StatefulSet/web-sts")
			Expect(cachedResource).Should(BeNil())
		})
	})
	Context("Namespace tests", Ordered, func() {
		BeforeAll(func() {
			klog.Info("Creating resources in a labelled namespace...")
			cmd := exec.Command("kubectl", "apply", "--wait=true", "-Rf", "manifests/labeled-ns/")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session.Out, 5).Should(gbytes.Say(".*deployment.apps.* [created|unchanged|configured].*"))
			Eventually(session.Out, 5).Should(gbytes.Say(".*statefulset.apps.* [created|unchanged|configured].*"))

			cmd = exec.Command("kubectl", "wait", "--for=condition=ready", "--timeout=180s", "pod", "-l", "parent=sts-labelled-ns", "-n", "labelled-ns")
			stsSession, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(stsSession.Out, 10).Should(gbytes.Say(".*pod.* met*"))

			cmd = exec.Command("kubectl", "wait", "--for=condition=ready", "--timeout=180s", "pod", "-l", "parent=deploy-labelled-ns", "-n", "labelled-ns")
			deploySession, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(deploySession.Out, 10).Should(gbytes.Say(".*pod.* met*"))

			cmd = exec.Command("kubectl", "wait", "--for=condition=ready", "--timeout=180s", "pod", "-l", "parent=deploy-excluded", "-n", "labelled-ns")
			excludedSession, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(excludedSession.Out, 10).Should(gbytes.Say(".*pod.* met*"))

			//re-initialize cache and confirm content
			controllers.Cache.Flush()
			err = controllers.InitializeCache()
			if err != nil {
				klog.Fatalf("Error rebuilding cache: %s", err.Error())
			}
		})
		AfterAll(func() {
			klog.Info("Cleaning up resources in a labelled namespace...")
			cmd := exec.Command("kubectl", "delete", "--wait=true", "-Rf", "manifests/labeled-ns/")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session.Out, 5).Should(gbytes.Say(".*deleted.*"))
		})
		It("Can fetch StatefulSets in a labelled namespace", func() {
			Eventually(session.Err, 20).Should(gbytes.Say(`.*[Added|Updated] StatefulSet: "labelled-ns-web-sts"*`))

			cachedResource, _ := controllers.Cache.Get("ResourceType/labelled-ns/StatefulSet/labelled-ns-web-sts")
			Expect(cachedResource).ShouldNot(BeNil())
		})
		It("Can fetch Deployments in a labelled namespace", func() {
			Eventually(session.Err, 20).Should(gbytes.Say(`.*[Added|Updated] Deployment: "labelled-ns-web-deploy"*`))

			cachedResource, _ := controllers.Cache.Get("ResourceType/labelled-ns/Deployment/labelled-ns-web-deploy")
			Expect(cachedResource).ShouldNot(BeNil())
		})
		It("Can not find Excluded resources", func() {
			Eventually(session.Err, 20).ShouldNot(gbytes.Say(`.*[Added|Updated] Deployment: "labelled-ns-disbaled-web-deploy"*`))

			cachedResource, _ := controllers.Cache.Get("ResourceType/labelled-ns/StatefulSet/labelled-ns-disbaled-web-deploy")
			Expect(cachedResource).Should(BeNil())
		})
	})
})
