package ardoq_test

import (
	"K8SArdoqBridge/app/controllers"
	"K8SArdoqBridge/app/tests/helper"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"

	_ "github.com/go-task/slim-sprig"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	_ "golang.org/x/tools/go/ast/inspector"
)

func TestArdoqController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ardoq Controller Suite")
}

var (
	tempClusterName = helper.RandomString(5) + "-adq-test"
	mockServer      *helper.MockArdoqServer
)

var _ = BeforeSuite(func() {
	log.Info("Initializing Mock Ardoq Server")

	// Create and configure mock server
	mockServer = helper.NewMockArdoqServer()
	err := mockServer.ConfigureMockEnvironment()
	Expect(err).NotTo(HaveOccurred())

	err = os.Setenv("ARDOQ_CLUSTER", tempClusterName)
	if err != nil {
		log.Error(err)
	}

	log.Infof("Mock server URL: %s", mockServer.URL)
	log.Infof("Cluster to be used: %s", tempClusterName)

	// Initialize cache
	controllers.Cache.Flush()
	err = controllers.InitializeCache()
	if err != nil {
		log.Fatalf("Error building cache: %s", err.Error())
	}

	// Note: BootstrapModel and BootstrapFields are skipped since the mock server
	// already has the model configured. If your tests need these, you'll need to
	// add mock endpoints for model/field creation.

	log.Info("Initializing Complete")
})

var _ = AfterSuite(func() {
	//_ = controllers.GenericDelete("Cluster", tempClusterName)
	if mockServer != nil {
		mockServer.Close()
	}
	controllers.Cache.Flush()
	log.Info("Cleanup Complete...Terminating!!")
})
