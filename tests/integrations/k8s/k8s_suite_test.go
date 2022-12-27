package k8s_test

import (
	"K8SArdoqBridge/app/controllers"
	"K8SArdoqBridge/app/tests/helper"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/util/homedir"
	"os"
	"os/exec"
	"path/filepath"
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
	log.Info("Initializing")
	err := os.Setenv("ARDOQ_CLUSTER", helper.RandomString(5)+"-k8s-"+os.Getenv("ARDOQ_CLUSTER"))
	if err != nil {
		log.Error(err)
	}
	log.Infof("Creating cluster: %s", os.Getenv("ARDOQ_CLUSTER"))

	publisherPath, err := gexec.Build("../../../main.go")
	Expect(err).NotTo(HaveOccurred())
	var cmd *exec.Cmd
	if os.Getenv("KUBECONFIG") != "" {
		cmd = exec.Command(publisherPath, "--kubeconfig", os.Getenv("KUBECONFIG"))
	} else if home := homedir.HomeDir(); home != "" {
		cmd = exec.Command(publisherPath, "--kubeconfig", filepath.Join(home, ".kube", "config"))
	}
	session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session.Err, 5).Should(gbytes.Say(".*Got watcher client.*"))
	Eventually(session.Err, 30).Should(gbytes.Say(`.*Initialised cluster in Ardoq`))
	Eventually(session.Err, 10).Should(gbytes.Say(`.*Starting event buffer`))
	Eventually(session.Err, 30).Should(gbytes.Say(`.*successfully acquired lease.*`))
	controllers.ApplyDelay(5)
	log.Info("Initializing Complete")
})

var _ = AfterSuite(func() {
	log.Info("Cleanup")
	//Cleanup shared components
	helper.CleanupSharedComponents("Resource")
	helper.CleanupSharedComponents("Node")

	//cleanup cluster in ardoq
	cleanupCluster()
	//kill the running session
	session.Kill()
	//cleanup running binary
	gexec.CleanupBuildArtifacts()
	log.Info("Cleanup Complete...Terminating!!")
})

func cleanupCluster() {
	resp, err := controllers.RestyClient().
		SetQueryParams(map[string]string{
			"workspace": os.Getenv("ARDOQ_WORKSPACE_ID"),
			"name":      os.Getenv("ARDOQ_CLUSTER")}).
		Get("component/search")
	if err != nil {
		log.Errorf("Error fetching clusters: %s", err)
	}
	var cluster []controllers.Component
	_ = controllers.Decode(resp.Body(), &cluster)
	for _, v := range cluster {
		_, err = controllers.RestyClient().Delete("component/" + v.ID)
		if err != nil {
			log.Errorf("Error deleting cluster %s: %s", v.Name, err)

		}
	}
}
