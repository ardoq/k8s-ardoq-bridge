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
func (b *BridgeController) ControlLoop(cancelContext context.Context) {

	for {
		select {
		case <-cancelContext.Done():
			break
		}
	}
}
