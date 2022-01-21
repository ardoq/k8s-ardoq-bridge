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
		It("Can fetch created Deployments", func() {
			Eventually(func() float64 {
				data, err := controllers.ApplicationResourceSearch("default", "Deployment", "web-deploy")
				Expect(err).ShouldNot(HaveOccurred())
				parsedData := data.Path("total").Data().(float64)
				return parsedData
			}, 20).ShouldNot(BeZero())
		})
		It("Can not find deleted Deployments", func() {
			klog.Info("Deleting deployment...")
			cmd := exec.Command("kubectl", "delete", "--wait=true", "-f", "manifests/deployment.yaml")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session.Out, 5).Should(gbytes.Say(".*deployment.apps.* deleted.*"))

			cmd = exec.Command("kubectl", "wait", "--for=delete", "--timeout=180s", "pod", "-l", "app=nginx,parent=deploy")
			session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			session.ExitCode()
			Eventually(session.ExitCode, 10).Should(BeZero())
			klog.Info("Deleted deployment.")

			Eventually(func() float64 {
				data, err := controllers.ApplicationResourceSearch("default", "Deployment", "web-deploy")
				Expect(err).ShouldNot(HaveOccurred())
				parsedData := data.Path("total").Data().(float64)
				return parsedData
			}, 20).Should(BeZero())
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
			Eventually(func() float64 {
				data, err := controllers.ApplicationResourceSearch("default", "StatefulSet", "web-sts")
				Expect(err).ShouldNot(HaveOccurred())
				parsedData := data.Path("total").Data().(float64)
				return parsedData
			}, 20).ShouldNot(BeZero())
		})
		It("Can not find deleted StatefulSets", func() {
			klog.Info("Deleting statefulset...")
			cmd := exec.Command("kubectl", "delete", "--wait=true", "-f", "manifests/statefulset.yaml")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session.Out, 5).Should(gbytes.Say(".*statefulset.apps.* deleted.*"))

			cmd = exec.Command("kubectl", "wait", "--for=delete", "--timeout=180s", "pod", "-l", "app=nginx,parent=sts")
			session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			session.ExitCode()
			Eventually(session.ExitCode, 10).Should(BeZero())
			klog.Info("Deleted statefulset.")

			Eventually(func() float64 {
				data, err := controllers.ApplicationResourceSearch("default", "StatefulSet", "web-sts")
				Expect(err).ShouldNot(HaveOccurred())
				parsedData := data.Path("total").Data().(float64)
				return parsedData
			}, 20).Should(BeZero())
		})
	})
})
