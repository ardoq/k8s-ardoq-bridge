package ardoq_test

import (
	"K8SArdoqBridge/app/controllers"
	"K8SArdoqBridge/app/tests/helper"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
	"time"
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
			ContainerRuntime:  "docker://20.10.11",
			KernelVersion:     "5.10.76-linuxkit",
			KubeletVersion:    "v1.22.4",
			KubeProxyVersion:  "v1.22.4",
			OperatingSystem:   "linux",
			OSImage:           "Docker Desktop",
			Provider:          "docker-desktop",
			CreationTimestamp: v1.Now().Format(time.RFC3339),
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
	Context("Node to Ardoq Integration tests", Ordered, func() {
		var compId string
		AfterAll(func() {
			_ = controllers.GenericDeleteSharedComponents("Node", "node_os", node.OperatingSystem)
			_ = controllers.GenericDeleteSharedComponents("Node", "architecture", node.Architecture)
			log.Info("Cleaned up shared node components")
		})
		It("Can create Node", func() {
			compId = controllers.GenericUpsert("Node", *node)
			Expect(compId).ShouldNot(BeNil())
		})
		It("Shared node components created", func() {
			controllers.Cache.Flush()
			err := controllers.InitializeCache()
			if err != nil {
				log.Fatalf("Error rebuilding cache: %s", err.Error())
			}
			cachedResource, found := controllers.Cache.Get("SharedNodeComponent/node_os/" + strings.ToLower(node.OperatingSystem))
			Expect(cachedResource).ShouldNot(BeNil())
			Expect(found).Should(BeTrue())

			cachedResource, found = controllers.Cache.Get("SharedNodeComponent/architecture/" + strings.ToLower(node.Architecture))
			Expect(cachedResource).ShouldNot(BeNil())
			Expect(found).Should(BeTrue())
		})
		It("Shared node components links created", func() {
			controllers.Cache.Flush()
			err := controllers.InitializeCache()
			if err != nil {
				log.Fatalf("Error rebuilding cache: %s", err.Error())
			}
			cachedResource, found := controllers.Cache.Get("SharedNodeLinks/" + compId + "/" + controllers.GenericUpsertSharedComponents("Node", "node_os", node.OperatingSystem))
			Expect(cachedResource).ShouldNot(BeNil())
			Expect(found).Should(BeTrue())

			cachedResource, found = controllers.Cache.Get("SharedNodeLinks/" + compId + "/" + controllers.GenericUpsertSharedComponents("Node", "architecture", node.Architecture))
			Expect(cachedResource).ShouldNot(BeNil())
			Expect(found).Should(BeTrue())
		})
		It("Can Update Node", func() {
			node.KernelVersion = "5.10.76-linuxkit-2"
			node.KubeletVersion = "v1.23.0"
			node.KubeProxyVersion = "v1.23.0"
			Expect(controllers.GenericUpsert("Node", *node)).ShouldNot(BeNil())
		})
		It("Can Delete Node", func() {
			Expect(controllers.GenericDelete("Node", *node)).Should(BeNil())
		})
		It("Can't Delete Non Existent Node", func() {
			node.Name = helper.RandomString(10)
			Expect(controllers.GenericDelete("Node", *node)).ShouldNot(BeNil())
		})
	})
})
