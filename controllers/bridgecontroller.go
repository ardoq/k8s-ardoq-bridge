package controllers

import (
	"context"
	v1 "k8s.io/api/apps/v1"
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"
	"strings"
)

type BridgeController struct {
	KubeClient *kubernetes.Clientset
}

func GetContainerImages(containers []v12.Container) string {
	values := make([]string, 0, len(containers))
	for _, v := range containers {
		values = append(values, v.Image)
	}
	return strings.Join(values, ", ")
}

func (b *BridgeController) OnDeploymentEvent(event watch.Event, res *v1.Deployment) {
	resourceType := "Deployment"
	//klog.Infof("%s | Resource: %s | %s | %s | %d | %s ", event.Type, res.Namespace, resourceType, res.Name, res.Status.Replicas, res.Spec.Template.Spec.Containers[0].Image)

	if res.Name == "" {
		klog.Errorf("Unable to retrieve %s from incoming event", resourceType)
		return
	}
	resource := Resource{
		RType:     resourceType,
		Name:      res.Name,
		ID:        "",
		Namespace: res.Namespace,
		Replicas:  int64(res.Status.Replicas),
		Image:     GetContainerImages(res.Spec.Template.Spec.Containers),
	}
	switch event.Type {
	case watch.Added, watch.Modified:
		UpsertDeploymentStatefulset(resource)
		break
	case watch.Deleted:
		DeleteDeploymentStatefulset(resource)
		break
	}
}
func (b *BridgeController) OnStatefulsetEvent(event watch.Event, res *v1.StatefulSet) {
	resourceType := "StatefulSet"
	//klog.Infof("%s | Resource: %s | %s | %s | %d | %s ", event.Type, res.Namespace, resourceType, res.Name, res.Status.Replicas, res.Spec.Template.Spec.Containers[0].Image)

	if res.Name == "" {
		klog.Errorf("Unable to retrieve %s from incoming event", resourceType)
		return
	}
	resource := Resource{
		RType:     resourceType,
		Name:      res.Name,
		ID:        "",
		Namespace: res.Namespace,
		Replicas:  int64(res.Status.Replicas),
		Image:     GetContainerImages(res.Spec.Template.Spec.Containers),
	}
	switch event.Type {
	case watch.Added, watch.Modified:
		UpsertDeploymentStatefulset(resource)
		break
	case watch.Deleted:
		DeleteDeploymentStatefulset(resource)
		break
	}
}
func (b *BridgeController) OnNodeEvent(event watch.Event, res *v12.Node) {
	resourceType := "Node"
	//klog.Infof("%s | Resource: %s  ", event.Type, res.ClusterName, res.Status.Capacity, res.Status.Allocatable, res.Status.NodeInfo.Architecture, res.Status.NodeInfo.ContainerRuntimeVersion, res.Status.NodeInfo.KernelVersion, res.Status.NodeInfo.KubeletVersion, res.Status.NodeInfo.KubeProxyVersion, res.Status.NodeInfo.OperatingSystem, res.Status.NodeInfo.OSImage)
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
			Storage: res.Status.Capacity.Storage().String(),
			Pods:    res.Status.Capacity.Pods().Value(),
		},
		Allocatable: NodeResources{
			CPU:     res.Status.Allocatable.Cpu().Value(),
			Memory:  res.Status.Allocatable.Memory().String(),
			Storage: res.Status.Allocatable.Storage().String(),
			Pods:    res.Status.Allocatable.Pods().Value(),
		},
		ContainerRuntime: res.Status.NodeInfo.ContainerRuntimeVersion,
		KernelVersion:    res.Status.NodeInfo.KernelVersion,
		KubeletVersion:   res.Status.NodeInfo.KubeletVersion,
		KubeProxyVersion: res.Status.NodeInfo.KubeProxyVersion,
		OperatingSystem:  res.Status.NodeInfo.OperatingSystem,
		OSImage:          res.Status.NodeInfo.OSImage,
	}
	switch event.Type {
	case watch.Added, watch.Modified:
		UpsertNode(node)
		break
	case watch.Deleted:
		DeleteNode(node)
		break
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
