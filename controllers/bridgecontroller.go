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

func (b *BridgeController) OnResourceEvent(event watch.Event) {

}
func NewBridgeController(kubeClient *kubernetes.Clientset) *BridgeController {

	return &BridgeController{
		KubeClient: kubeClient,
	}
}
func (b *BridgeController) OnDeploymentEvent(event watch.Event, res *v1.Deployment) {
	klog.Infof("%s | Resource: %s | Deployment | %s | %d | %s ", event.Type, res.Namespace, res.Name, res.Status.Replicas, res.Spec.Template.Spec.Containers[0].Image)
}
func (b *BridgeController) OnStatefulsetEvent(event watch.Event, res *v1.StatefulSet) {
	klog.Infof("%s | Resource: %s | StatefulSet | %s | %d | %s ", event.Type, res.Namespace, res.Name, res.Status.Replicas, res.Spec.Template.Spec.Containers[0].Image)
}
func (b *BridgeController) ControlLoop(cancelContext context.Context) {

	for {
		select {
		case <-cancelContext.Done():
			break
		}
	}
}
