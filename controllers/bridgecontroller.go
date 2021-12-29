package controllers

import (
	"context"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"
)

type BridgeController struct {
	KubeClient *kubernetes.Clientset
}

func NewBridgeController(kubeClient *kubernetes.Clientset) *BridgeController {
	return &BridgeController{
		KubeClient: kubeClient,
	}
}
func (b *BridgeController) OnDeploymentEvent(event watch.Event, res *v1.Deployment) {
	resourceType := "Deployment"
	klog.Infof("%s | Resource: %s | %s | %s | %d | %s ", event.Type, res.Namespace, resourceType, res.Name, res.Status.Replicas, res.Spec.Template.Spec.Containers[0].Image)

	if res.Name == "" {
		klog.Errorf("Unable to retrieve %s from incoming event", resourceType)
		return
	}
	switch event.Type {
	case watch.Added, watch.Modified:
		firstOrCreateDeploymentStatefulset(res.Name, resourceType, res.Namespace)
		break
	case watch.Deleted:
		//todo: add delete processing

		break
	}
}
func (b *BridgeController) OnStatefulsetEvent(event watch.Event, res *v1.StatefulSet) {
	resourceType := "StatefulSet"
	klog.Infof("%s | Resource: %s | %s | %s | %d | %s ", event.Type, res.Namespace, resourceType, res.Name, res.Status.Replicas, res.Spec.Template.Spec.Containers[0].Image)

	if res.Name == "" {
		klog.Errorf("Unable to retrieve %s from incoming event", resourceType)
		return
	}
	switch event.Type {
	case watch.Added, watch.Modified:
		firstOrCreateDeploymentStatefulset(res.Name, resourceType, res.Namespace)
		break
	case watch.Deleted:
		//todo: add delete processing
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
