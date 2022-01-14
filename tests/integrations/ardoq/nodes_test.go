package ardoq_test

import (
	"K8SArdoqBridge/app/controllers"
	"K8SArdoqBridge/app/tests/helper"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Nodes", Ordered, func() {
	var node *controllers.Node
	var resourceName = helper.RandomString(10)
	BeforeEach(func() {
		node = &controllers.Node{
			Name:         resourceName,
			Architecture: "amd64",
			Capacity: controllers.NodeResources{
				CPU:     4,
				Memory:  "8150980Ki",
				Storage: "61255492Ki",
				Pods:    110,
			},
			Allocatable: controllers.NodeResources{
				CPU:     4,
				Memory:  "8048580Ki",
				Storage: "56453061Ki",
				Pods:    110,
			},
			ContainerRuntime: "docker://20.10.11",
			KernelVersion:    "5.10.76-linuxkit",
			KubeletVersion:   "v1.22.4",
			KubeProxyVersion: "v1.22.4",
			OperatingSystem:  "linux",
			OSImage:          "Docker Desktop",
			Provider:         "docker-desktop",
		}

	})
	Context("Structural Validation", func() {
		It("Ensure Node Structure is valid", func() {
			Expect(node.IsNodeValid()).To(BeTrue())
		})
		It("Node name not set", func() {
			node.Name = ""
			Expect(node.IsNodeValid()).To(BeFalse())
		})
		It("Node architecture not set", func() {
			node.Architecture = ""
			Expect(node.IsNodeValid()).To(BeFalse())
		})
		It("Node container runtime not set", func() {
			node.ContainerRuntime = ""
			Expect(node.IsNodeValid()).To(BeFalse())
		})
		It("Invalid kernel version Type set", func() {
			node.KernelVersion = ""
			Expect(node.IsNodeValid()).To(BeFalse())
		})
		It("Node kubelet not set", func() {
			node.KubeletVersion = ""
			Expect(node.IsNodeValid()).To(BeFalse())
		})
		It("Node kube proxy not set", func() {
			node.KubeProxyVersion = ""
			Expect(node.IsNodeValid()).To(BeFalse())
		})
		It("Node os not set", func() {
			node.OperatingSystem = ""
			Expect(node.IsNodeValid()).To(BeFalse())
		})
		It("Node os image not set", func() {
			node.OSImage = ""
			Expect(node.IsNodeValid()).To(BeFalse())
		})
		It("Node allocatable expected type", func() {
			Expect(node.Allocatable).To(BeAssignableToTypeOf(controllers.NodeResources{}))
		})
		It("Node capacity expected type", func() {
			Expect(node.Capacity).To(BeAssignableToTypeOf(controllers.NodeResources{}))
		})

	})
	Context("Node Ardoq Link tests", Ordered, func() {
		It("Can create Node", func() {
			Expect(controllers.UpsertNode(*node)).ShouldNot(BeNil())
			helper.ApplyDelay()
		})
		It("Can Update Node", func() {
			node.KernelVersion = "5.10.76-linuxkit-2"
			node.KubeletVersion = "v1.23.0"
			node.KubeProxyVersion = "v1.23.0"
			Expect(controllers.UpsertNode(*node)).ShouldNot(BeNil())
			helper.ApplyDelay()
		})
		It("Can Delete Node", func() {
			Expect(controllers.DeleteNode(*node)).Should(BeNil())
		})
		It("Can't Delete None Existent Node", func() {
			node.Name = helper.RandomString(10)
			Expect(controllers.DeleteNode(*node)).ShouldNot(BeNil())
		})
	})
})
