package ardoq_test

import (
	"K8SArdoqBridge/app/controllers"
	"k8s.io/klog/v2"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestArdoqController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ardoq Controller Suite")
}

var _ = BeforeSuite(func() {

})

var _ = AfterSuite(func() {
	_ = controllers.GenericDelete("Cluster", os.Getenv("ARDOQ_CLUSTER"))
	klog.Info("Cleanup Complete...Terminating!!")
})
