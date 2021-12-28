package subscriptions

import (
	"ArdoqK8sBridge/app/controllers"
	"ArdoqK8sBridge/app/lib/subscription"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type DeploymentSubscriber struct {
	BridgeDataProvider *controllers.BridgeController
}

func (DeploymentSubscriber) WithElectedResource() interface{} {

	return &v1.Deployment{}
}

func (DeploymentSubscriber) WithEventType() []watch.EventType {

	return []watch.EventType{watch.Added, watch.Deleted, watch.Modified}
}

func (d DeploymentSubscriber) OnEvent(msg subscription.Message) {

	res := msg.Event.Object.(*v1.Deployment)
	if res.Labels["sync-to-ardoq"] != "" {
		d.BridgeDataProvider.OnDeploymentEvent(msg.Event, res)
	}
}
