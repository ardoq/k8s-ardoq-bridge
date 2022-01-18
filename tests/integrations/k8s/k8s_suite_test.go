package k8s_test

import (
	"K8SArdoqBridge/app/tests/helper"
	"k8s.io/klog"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestK8s(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "K8s Suite")
}

var _ = BeforeSuite(func() {
	klog.Info("Waiting for service to load")
	helper.ApplyDelay(15)
	klog.Info("Starting tests")
})
