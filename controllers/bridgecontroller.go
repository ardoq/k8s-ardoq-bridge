package controllers

import (
	"context"
	"reflect"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/apps/v1"
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type resourceQueue struct {
	action string
	data   Resource
}
type nodeQueue struct {
	action string
	data   Node
}

var (
	resourceQueues = make(chan resourceQueue)
	nodeQueues     = make(chan nodeQueue)
)

type BridgeController struct {
	KubeClient *kubernetes.Clientset
}

func (b *BridgeController) OnApplicationResourceEvent(event watch.Event, genericResource interface{}) {
	resourceType := reflect.TypeOf(genericResource).String()
	resource := Resource{}
	if strings.HasSuffix(resourceType, "Deployment") {
		res := genericResource.(v1.Deployment)
		if res.Name == "" {
			log.Errorf("Unable to retrieve %s from incoming event", resourceType)
			return
		}
		resource = Resource{
			Name:              res.Name,
			ResourceType:      "Deployment",
			Namespace:         res.Namespace,
			Replicas:          res.Status.Replicas,
			Image:             GetContainerImages(res.Spec.Template.Spec.Containers),
			CreationTimestamp: res.CreationTimestamp.Format(time.RFC3339),
			Stack:             res.Labels["ardoq/stack"],
			Team:              res.Labels["ardoq/team"],
			Project:           res.Labels["ardoq/project"],
			Requests:          GetAppResourceRequirements(res.Spec.Template.Spec.Containers, "requests"),
			Limits:            GetAppResourceRequirements(res.Spec.Template.Spec.Containers, "limits"),
		}
	} else if strings.HasSuffix(resourceType, "StatefulSet") {
		res := genericResource.(v1.StatefulSet)
		if res.Name == "" {
			log.Errorf("Unable to retrieve %s from incoming event", resourceType)
			return
		}
		resource = Resource{
			Name:              res.Name,
			ResourceType:      "StatefulSet",
			Namespace:         res.Namespace,
			Replicas:          res.Status.Replicas,
			Image:             GetContainerImages(res.Spec.Template.Spec.Containers),
			CreationTimestamp: res.CreationTimestamp.Format(time.RFC3339),
			Stack:             res.Labels["ardoq/stack"],
			Team:              res.Labels["ardoq/team"],
			Project:           res.Labels["ardoq/project"],
			Requests:          GetAppResourceRequirements(res.Spec.Template.Spec.Containers, "requests"),
			Limits:            GetAppResourceRequirements(res.Spec.Template.Spec.Containers, "limits"),
		}
	} else {
		log.Errorf("Invalid type: %s", resourceType)
		return
	}
	switch event.Type {
	case watch.Added, watch.Modified:
		toQueue := resourceQueue{
			action: "UPSERT",
			data:   resource,
		}
		resourceQueues <- toQueue
		break
	case watch.Deleted:
		toQueue := resourceQueue{
			action: "DELETE",
			data:   resource,
		}
		resourceQueues <- toQueue
		break
	}
}

func (b *BridgeController) OnNodeEvent(event watch.Event, res *v12.Node) {
	resourceType := "Node"
	if res.Name == "" {
		log.Errorf("Unable to retrieve %s from incoming event", resourceType)
		return
	}
	node := Node{
		Name:         res.Name,
		Architecture: res.Status.NodeInfo.Architecture,
		Capacity: NodeResources{
			CPU:     res.Status.Capacity.Cpu().Value(),
			Memory:  ParseToMB(res.Status.Capacity.Memory().Value()),
			Storage: ParseToMB(res.Status.Capacity.StorageEphemeral().Value()),
			Pods:    res.Status.Capacity.Pods().Value(),
		},
		Allocatable: NodeResources{
			CPU:     res.Status.Allocatable.Cpu().Value(),
			Memory:  ParseToMB(res.Status.Allocatable.Memory().Value()),
			Storage: ParseToMB(res.Status.Allocatable.StorageEphemeral().Value()),
			Pods:    res.Status.Allocatable.Pods().Value(),
		},
		ContainerRuntime:  res.Status.NodeInfo.ContainerRuntimeVersion,
		InstanceType:      res.Labels["node.kubernetes.io/instance-type"],
		KernelVersion:     res.Status.NodeInfo.KernelVersion,
		KubeletVersion:    res.Status.NodeInfo.KubeletVersion,
		KubeProxyVersion:  res.Status.NodeInfo.KubeProxyVersion,
		OperatingSystem:   res.Status.NodeInfo.OperatingSystem,
		OSImage:           res.Status.NodeInfo.OSImage,
		Provider:          res.Spec.ProviderID,
		Pool:              GetNodePool(res.Labels),
		CreationTimestamp: res.CreationTimestamp.Format(time.RFC3339),
		Region:            res.Labels["failure-domain.beta.kubernetes.io/region"],
		Zone:              res.Labels["failure-domain.beta.kubernetes.io/zone"],
	}
	switch event.Type {
	case watch.Added, watch.Modified:
		toQueue := nodeQueue{
			action: "UPSERT",
			data:   node,
		}
		nodeQueues <- toQueue
		break
	case watch.Deleted:
		toQueue := nodeQueue{
			action: "DELETE",
			data:   node,
		}
		nodeQueues <- toQueue
		break
	}
}

func ResourceConsumer() {
	for q := range resourceQueues {
		switch q.action {
		case "UPSERT":
			GenericUpsert(q.data.ResourceType, q.data)
			break
		case "DELETE":
			err := GenericDelete(q.data.ResourceType, q.data)
			if err != nil {
				return
			}
			break
		}

	}
}
func NodeConsumer() {
	for q := range nodeQueues {
		switch q.action {
		case "UPSERT":
			GenericUpsert("Node", q.data)
			break
		case "DELETE":
			err := GenericDelete("Node", q.data)
			if err != nil {
				return
			}
			break
		}

	}
}

func (b *BridgeController) ControlLoop(cancelContext context.Context) {

	for {
		select {
		case <-cancelContext.Done():
			break
		}
	}
}
