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
	if PerformSync(getNamespaceLabels(res.Namespace), res.Labels) {
		d.BridgeDataProvider.OnApplicationResourceEvent(msg.Event, *res)
	}

	//Perform cleanup if it was previously labeled, or we are booting /Reconciliation
	if (msg.Event.Type == watch.Modified || msg.Event.Type == watch.Added) && !PerformSync(getNamespaceLabels(res.Namespace), res.Labels) {
		resource := controllers.Resource{
			ResourceType: "Deployment",
			Name:         res.Name,
			ID:           "",
			Namespace:    res.Namespace,
			Replicas:     res.Status.Replicas,
			Image:        controllers.GetContainerImages(res.Spec.Template.Spec.Containers),
		}
		err := controllers.GenericDelete(resource.ResourceType, resource)
		if err != nil {
			return
		}
	}
}
