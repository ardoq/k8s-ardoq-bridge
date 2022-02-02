package controllers

import (
	"context"
	v1 "k8s.io/api/apps/v1"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	"reflect"
	"strings"
	"time"
)

var (
	upsertQueue = make(chan Resource)
	deleteQueue = make(chan Resource)
)

type BridgeController struct {
	KubeClient *kubernetes.Clientset
}

func GetContainerImages(containers []v12.Container) string {
	values := make([]string, 0, len(containers))
	for _, v := range containers {
		values = append(values, v.Image)
	}
	return strings.Join(values, ",")
}
func (b *BridgeController) OnApplicationResourceEvent(event watch.Event, genericResource interface{}) {
	resourceType := reflect.TypeOf(genericResource).String()
	resource := Resource{}
	if strings.HasSuffix(resourceType, "Deployment") {
		res := genericResource.(v1.Deployment)
		if res.Name == "" {
			klog.Errorf("Unable to retrieve %s from incoming event", resourceType)
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
		}
	} else if strings.HasSuffix(resourceType, "StatefulSet") {
		res := genericResource.(v1.StatefulSet)
		if res.Name == "" {
			klog.Errorf("Unable to retrieve %s from incoming event", resourceType)
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
		}
	} else {
		klog.Errorf("Invalid type: %s", resourceType)
		return
	}
	switch event.Type {
	case watch.Added, watch.Modified:
		upsertQueue <- resource
		break
	case watch.Deleted:
		deleteQueue <- resource
		break
	}
}
func (b *BridgeController) OnNamespaceEvent(event watch.Event, res *v12.Namespace) {
	resourceType := "Namespace"
	if res.Name == "" {
		klog.Errorf("Unable to retrieve %s from incoming event", resourceType)
		return
	}

	deploys, err := ClientSet.AppsV1().Deployments(res.Name).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		klog.Errorf("failed to get List deployments: %w", err)
	}
	for _, d := range deploys.Items {
		if d.GetLabels()["sync-to-ardoq"] != "disabled" {
			b.OnApplicationResourceEvent(event, d)
		}

	}
	sts, err := ClientSet.AppsV1().StatefulSets(res.Name).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		klog.Errorf("failed to get List deployments: %v", err)
	}
	for _, d := range sts.Items {
		if d.GetLabels()["sync-to-ardoq"] != "disabled" {
			b.OnApplicationResourceEvent(event, d)
		}
	}
}
func (b *BridgeController) OnNodeEvent(event watch.Event, res *v12.Node) {
	resourceType := "Node"
	if res.Name == "" {
		klog.Errorf("Unable to retrieve %s from incoming event", resourceType)
		return
	}
	node := Node{
		Name:         res.Name,
		Architecture: res.Status.NodeInfo.Architecture,
		Capacity: NodeResources{
			CPU:     res.Status.Capacity.Cpu().Value(),
			Memory:  res.Status.Capacity.Memory().String(),
			Storage: res.Status.Capacity.StorageEphemeral().String(),
			Pods:    res.Status.Capacity.Pods().Value(),
		},
		Allocatable: NodeResources{
			CPU:     res.Status.Allocatable.Cpu().Value(),
			Memory:  res.Status.Allocatable.Memory().String(),
			Storage: res.Status.Allocatable.StorageEphemeral().String(),
			Pods:    res.Status.Allocatable.Pods().Value(),
		},
		ContainerRuntime:  res.Status.NodeInfo.ContainerRuntimeVersion,
		KernelVersion:     res.Status.NodeInfo.KernelVersion,
		KubeletVersion:    res.Status.NodeInfo.KubeletVersion,
		KubeProxyVersion:  res.Status.NodeInfo.KubeProxyVersion,
		OperatingSystem:   res.Status.NodeInfo.OperatingSystem,
		OSImage:           res.Status.NodeInfo.OSImage,
		Provider:          res.Spec.ProviderID,
		CreationTimestamp: res.CreationTimestamp.Format(time.RFC3339),
	}
	switch event.Type {
	case watch.Added, watch.Modified:
		GenericUpsert("Node", node)
		break
	case watch.Deleted:
		err := GenericDelete("Node", node)
		if err != nil {
			return
		}
		break
	}
}

func ResourceUpsertConsumer() {
	for res := range upsertQueue {
		GenericUpsert(res.ResourceType, res)
	}
}
func ResourceDeleteConsumer() {
	for res := range deleteQueue {
		err := GenericDelete(res.ResourceType, res)
		if err != nil {
			return
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
