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

var tempClusterName = helper.RandomString(5) + "-adq-" + os.Getenv("ARDOQ_CLUSTER")

var _ = BeforeSuite(func() {
	log.Info("Initializing")
	err := os.Setenv("ARDOQ_CLUSTER", tempClusterName)
	if err != nil {
		log.Error(err)
	}
	log.Infof("Cluster to be used: %s", tempClusterName)
	controllers.Cache.Flush()
	err = controllers.InitializeCache()
	if err != nil {
		log.Fatalf("Error building cache: %s", err.Error())
	}
	//Initialise the Model in Ardoq
	err = controllers.BootstrapModel()
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("Initialized the Model")

	//Initialise the Custom Fields
	err = controllers.BootstrapFields()
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("Initialised Custom Fields")
	log.Info("Initializing Complete")
})

var _ = AfterSuite(func() {
	_ = controllers.GenericDelete("Cluster", tempClusterName)
	controllers.Cache.Flush()
	log.Info("Cleanup Complete...Terminating!!")
})
