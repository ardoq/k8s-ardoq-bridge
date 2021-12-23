package subscriptions

import (
	"KubeOps/app/lib/subscription"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/klog"
)

type DeploymentOperator struct{}

func (DeploymentOperator) WithElectedResource() interface{} {

	return &v1.Deployment{}
}

func (DeploymentOperator) WithEventType() []watch.EventType {

	return []watch.EventType{watch.Added, watch.Deleted, watch.Modified}
}

func (DeploymentOperator) OnEvent(msg subscription.Message) {

	deploy := msg.Event.Object.(*v1.Deployment)
	if deploy.Labels["sync-to-ardoq"] != "" {
		klog.Info(msg.Event.Type)
		klog.Infof("Deployment %s has %d Available replicas", deploy.Name, deploy.Status.AvailableReplicas)
	}
}
