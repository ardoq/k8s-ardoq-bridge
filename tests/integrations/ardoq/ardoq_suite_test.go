package ardoq_test

import (
	"K8SArdoqBridge/app/controllers"
	"K8SArdoqBridge/app/tests/helper"
	"k8s.io/klog/v2"
	"os"
	"testing"

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
	klog.Info("Initializing")
	err := os.Setenv("ARDOQ_CLUSTER", tempClusterName)
	if err != nil {
		klog.Error(err)
	}
	klog.Infof("Cluster to be used: %q", tempClusterName)
	klog.Info("Initializing Complete")
})

var _ = AfterSuite(func() {
	_ = controllers.GenericDelete("Cluster", tempClusterName)
	klog.Info("Cleanup Complete...Terminating!!")
})
