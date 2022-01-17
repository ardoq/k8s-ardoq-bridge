package ardoq_test

import (
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

})
