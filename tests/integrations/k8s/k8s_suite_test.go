package k8s_test

import (
	"K8SArdoqBridge/app/controllers"
	"K8SArdoqBridge/app/tests/helper"
	"context"
	"fmt"
	ardoq "github.com/mories76/ardoq-client-go/pkg"
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
	Eventually(session.Err, 20).Should(gbytes.Say(`.*Initialised cluster in Ardoq`))
	Eventually(session.Err, 10).Should(gbytes.Say(`.*Starting event buffer`))
	Eventually(session.Err, 20).Should(gbytes.Say(`.*successfully acquired lease.*`))
	controllers.ApplyDelay(5)
	log.Info("Initializing Complete")
})

var _ = AfterSuite(func() {
	log.Info("Cleanup")
	//cleanup cluster in ardoq
	cleanupCluster()
	//cleanup running binary
	session.Kill()
	gexec.CleanupBuildArtifacts()
	log.Info("Cleanup Complete...Terminating!!")
})

func cleanupCluster() {
	a, err := ardoq.NewRestClient(os.Getenv("ARDOQ_BASEURI"), os.Getenv("ARDOQ_APIKEY"), os.Getenv("ARDOQ_ORG"), "v0.0.0")
	if err != nil {
		fmt.Printf("cannot create new restclient %s", err)
		os.Exit(1)
	}
	cluster, err := a.Components().Search(context.TODO(), &ardoq.ComponentSearchQuery{Workspace: os.Getenv("ARDOQ_WORKSPACE_ID"), Name: os.Getenv("ARDOQ_CLUSTER")})
	if err != nil {
		log.Errorf("Error fetching clusters: %s", err)
	}

	for _, v := range *cluster {
		err = a.Components().Delete(context.TODO(), v.ID)
		if err != nil {
			log.Errorf("Error deleting cluster %s: %s", v.Name, err)

		}
	}
}
