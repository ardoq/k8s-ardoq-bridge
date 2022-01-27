package subscriptions

import (
	"K8SArdoqBridge/app/controllers"
	"K8SArdoqBridge/app/lib/subscription"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type NamespaceSubscriber struct {
	BridgeDataProvider *controllers.BridgeController
}

func (NamespaceSubscriber) WithElectedResource() interface{} {

	return &v1.Namespace{}
}

func (NamespaceSubscriber) WithEventType() []watch.EventType {

	return []watch.EventType{watch.Added, watch.Deleted, watch.Modified}
}

func (d NamespaceSubscriber) OnEvent(msg subscription.Message) {
	res := msg.Event.Object.(*v1.Namespace)
	if res.Labels["sync-to-ardoq"] == "enabled" {
		d.BridgeDataProvider.OnNamespaceEvent(msg.Event, res)
	}
	//Perform cleanup if it was previously labeled, or we are booting
	if (msg.Event.Type == watch.Modified || msg.Event.Type == watch.Added) && (res.Labels["sync-to-ardoq"] == "" || res.Labels["sync-to-ardoq"] == "disabled") {

	}
}
