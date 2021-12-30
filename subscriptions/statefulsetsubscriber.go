package subscriptions

import (
	"ArdoqK8sBridge/app/controllers"
	"ArdoqK8sBridge/app/lib/subscription"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type StatefulsetSubscriber struct {
	BridgeDataProvider *controllers.BridgeController
}

func (StatefulsetSubscriber) WithElectedResource() interface{} {

	return &v1.StatefulSet{}
}

func (StatefulsetSubscriber) WithEventType() []watch.EventType {

	return []watch.EventType{watch.Added, watch.Deleted, watch.Modified}
}

func (d StatefulsetSubscriber) OnEvent(msg subscription.Message) {

	res := msg.Event.Object.(*v1.StatefulSet)
	if res.Labels["sync-to-ardoq"] != "" {
		d.BridgeDataProvider.OnStatefulsetEvent(msg.Event, res)
	}
	if msg.Event.Type == watch.Modified && (res.Labels["sync-to-ardoq"] == "") {
		resource := controllers.Resource{
			RType:     "StatefulSet",
			Name:      res.Name,
			ID:        "",
			Namespace: res.Namespace,
			Replicas:  int64(res.Status.Replicas),
			Image:     controllers.GetContainerImages(res.Spec.Template.Spec.Containers),
		}
		controllers.DeleteDeploymentStatefulset(resource)
	}
}
