package subscriptions

import (
	"ArdoqK8sBridge/app/controllers"
	"ArdoqK8sBridge/app/lib/subscription"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/klog"
)

type NodeSubscriber struct {
	BridgeDataProvider *controllers.BridgeController
}

func (NodeSubscriber) WithElectedResource() interface{} {

	return &v1.Node{}
}

func (NodeSubscriber) WithEventType() []watch.EventType {

	return []watch.EventType{watch.Added, watch.Deleted, watch.Modified}
}

func (d NodeSubscriber) OnEvent(msg subscription.Message) {
	res := msg.Event.Object.(*v1.Node)
	if res.Name == "" {
		klog.Errorf("Unable to retrieve Node from incoming event")
		return
	}
	d.BridgeDataProvider.OnNodeEvent(msg.Event, res)
}
