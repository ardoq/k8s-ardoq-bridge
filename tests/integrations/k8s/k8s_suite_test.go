package k8s_test

import (
	"K8SArdoqBridge/app/controllers"
	"K8SArdoqBridge/app/tests/helper"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestK8s(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "K8s Suite")
}

var (
	fakeK8sClient kubernetes.Interface
	mockServer    *helper.MockArdoqServer
	clusterName   string
)

var _ = BeforeSuite(func() {
	log.Info("Initializing Mock Ardoq Server")
	mockServer = helper.NewMockArdoqServer()
	err := mockServer.ConfigureMockEnvironment()
	Expect(err).NotTo(HaveOccurred())
	log.Infof("Mock server URL: %s", mockServer.URL)

	log.Info("Initializing Fake K8s Client")
	fakeK8sClient = helper.NewFakeK8sClient()
	controllers.ClientSet = fakeK8sClient

	clusterName = helper.RandomString(5) + "-k8s-test"
	err = os.Setenv("ARDOQ_CLUSTER", clusterName)
	Expect(err).NotTo(HaveOccurred())
	log.Infof("Cluster name: %s", clusterName)

	// Initialize Ardoq model and fields
	err = controllers.BootstrapModel()
	Expect(err).NotTo(HaveOccurred())
	log.Info("Initialized the Model")

	err = controllers.BootstrapFields()
	Expect(err).NotTo(HaveOccurred())
	log.Info("Initialised Custom Fields")

	// Initialize cache
	err = controllers.InitializeCache()
	Expect(err).NotTo(HaveOccurred())
	log.Info("Cache initialized")

	// Initialize cluster in Ardoq
	controllers.LookupCluster(clusterName)
	log.Info("Initialised cluster in Ardoq")

	log.Info("Initializing Complete")
})

var _ = AfterSuite(func() {
	log.Info("Cleanup")
	//Cleanup shared components
	helper.CleanupSharedComponents("Resource")
	helper.CleanupSharedComponents("Node")

	//cleanup cluster in ardoq
	cleanupCluster()

	// Close mock server
	if mockServer != nil {
		mockServer.Close()
	}

	controllers.Cache.Flush()
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
