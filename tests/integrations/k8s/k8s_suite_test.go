package k8s_test

import (
	"K8SArdoqBridge/app/tests/helper"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"k8s.io/klog"
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestK8s(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "K8s Suite")
}

var (
	session *gexec.Session
)

var _ = BeforeSuite(func() {
	klog.Info("Initializing")
	publisherPath, err := gexec.Build("../../../main.go")
	Expect(err).NotTo(HaveOccurred())
	cmd := exec.Command(publisherPath)
	session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session.Err, 5).Should(gbytes.Say(".*Got watcher client.*"))
	Eventually(session.Err, 5).Should(gbytes.Say(`.*Initialising cluster in Ardoq`))
	Eventually(session.Err, 5).Should(gbytes.Say(`.*Starting event buffer`))
	Eventually(session.Err, 15).Should(gbytes.Say(`.*successfully acquired lease.*`))
	helper.ApplyDelay(10)
	klog.Info("Initializing Complete")
})

var _ = AfterSuite(func() {
	klog.Info("Cleanup")
	session.Kill()
	gexec.CleanupBuildArtifacts()
	klog.Info("Cleanup Complete...Terminating!!")
})
