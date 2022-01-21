package ardoq_test

import (
	"K8SArdoqBridge/app/controllers"
	"K8SArdoqBridge/app/tests/helper"
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
	klog.Info("Initializing")
	err := os.Setenv("ARDOQ_CLUSTER", helper.RandomString(5)+"-"+os.Getenv("ARDOQ_CLUSTER"))
	if err != nil {
		klog.Error(err)
	}
	klog.Infof("Creating cluster: %q", os.Getenv("ARDOQ_CLUSTER"))
	klog.Info("Initializing Complete")
})

var _ = AfterSuite(func() {
	_ = controllers.GenericDelete("Cluster", os.Getenv("ARDOQ_CLUSTER"))
	klog.Info("Cleanup Complete...Terminating!!")
})
