package subscriptions

import (
	"K8SArdoqBridge/app/controllers"
	"K8SArdoqBridge/app/lib/subscription"
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
	if msg.Event.Type == watch.Modified && (res.Labels["sync-to-ardoq"] == "") {
		resource := controllers.Resource{
			ResourceType: "Deployment",
			Name:         res.Name,
			ID:           "",
			Namespace:    res.Namespace,
			Replicas:     int64(res.Status.Replicas),
			Image:        controllers.GetContainerImages(res.Spec.Template.Spec.Containers),
		}
		err := controllers.DeleteApplicationResource(resource)
		if err != nil {
			return
		}
	}
}
