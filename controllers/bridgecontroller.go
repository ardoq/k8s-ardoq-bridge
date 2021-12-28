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
	klog.Info(event.Type)
	klog.Infof("Deployment %s has %d Available replicas", res.Name, res.Status.AvailableReplicas)
}
func (b *BridgeController) ControlLoop(cancelContext context.Context) {

	for {
		select {
		case <-cancelContext.Done():
			break
		}
	}
}
