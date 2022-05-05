package controllers

import (
	"K8SArdoqBridge/app/lib/metrics"
	"context"
	"errors"
	goCache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	v12 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"strconv"
	"strings"
	"time"
)

var (
	Cache     = goCache.New(5*time.Minute, 10*time.Minute)
	ClientSet *kubernetes.Clientset
)

func LookupCluster(name string, deletion ...bool) string {
	if cachedResource, found := GetFromCache("ResourceType/Cluster/" + name); found {
		return cachedResource.(string)
	}
	if !(len(deletion) > 0 && deletion[0]) {
		return GenericUpsert("Cluster", name)
	}
	return ""
}
func LookupNamespace(name string) string {
	if cachedResource, found := GetFromCache("ResourceType/Namespace/" + name); found {
		return cachedResource.(string)
	}
	return ""
}

func LookupResource(namespace string, resourceType string, resourceName string) string {
	if cachedResource, found := GetFromCache("ResourceType/" + namespace + "/" + resourceType + "/" + resourceName); found {
		return cachedResource.(Resource).ID
	}
	return ""
}
func LookupNode(name string) string {
	if cachedResource, found := GetFromCache("ResourceType/Node/" + name); found {
		return cachedResource.(Node).ID
	}
	return ""
}

func lookUpTypeId(name string) string {
	if typeId, found := GetFromCache("ArdoqTypes/" + name); found {
		return typeId.(string)
	}
	requestStarted := time.Now()
	workspace, err := ardRestClient().Workspaces().Get(context.TODO(), workspaceId)
	metrics.RequestLatency.WithLabelValues("read").Observe(time.Since(requestStarted).Seconds())
	if err != nil {
		metrics.RequestStatusCode.WithLabelValues("error").Inc()
		log.Errorf("Error getting workspace: %s", err)
	}
	metrics.RequestStatusCode.WithLabelValues("success").Inc()
	//set componentModel to the componentModel from the found workspace
	componentModel := workspace.ComponentModel
	requestStarted = time.Now()
	model, err := ardRestClient().Models().Read(context.TODO(), componentModel)
	metrics.RequestLatency.WithLabelValues("read").Observe(time.Since(requestStarted).Seconds())
	if err != nil {
		metrics.RequestStatusCode.WithLabelValues("error").Inc()
		log.Errorf("Error getting model: %s", err)
	}
	metrics.RequestStatusCode.WithLabelValues("success").Inc()
	cmpTypes := model.GetComponentTypeID()
	if cmpTypes[name] != "" {
		PersistToCache("ArdoqTypes/"+name, cmpTypes[name])
		return cmpTypes[name]
	} else {
		return ""
	}

}

func (r *Resource) IsApplicationResourceValid() bool {
	if r.Name != "" && r.Namespace != "" && r.ResourceType != "" && r.Image != "" && Contains(validApplicationResourceTypes, r.ResourceType) {
		return true
	}
	return false
}
func (n *Node) IsNodeValid() bool {
	if n.Name != "" && n.Architecture != "" && n.KernelVersion != "" && n.KubeletVersion != "" && n.KubeProxyVersion != "" && n.OperatingSystem != "" && n.OSImage != "" && n.ContainerRuntime != "" {
		return true
	}
	return false
}
func ApplyDelay(seconds ...time.Duration) {
	if len(seconds) > 0 {
		time.Sleep(seconds[0] * time.Second)
	} else {
		time.Sleep(5 * time.Second)
	}
}
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
func GetFromCache(name string) (interface{}, bool) {
	if cachedResource, found := Cache.Get(name); found {
		metrics.CacheHits.Inc()
		return cachedResource, true
	} else {
		metrics.CacheMiss.Inc()
		return nil, false
	}
}
func PersistToCache(name string, value interface{}) {
	Cache.Set(name, value, goCache.NoExpiration)
	metrics.CachePersists.Inc()
}
func GetNodePool(nodeLabels map[string]string) string {
	if nodeLabels["kubernetes.azure.com/agentpool"] != "" {
		return nodeLabels["kubernetes.azure.com/agentpool"]
	} else if nodeLabels["eks.amazonaws.com/nodegroup"] != "" {
		return nodeLabels["eks.amazonaws.com/nodegroup"]
	}
	return ""
}
func GetContainerImages(containers []v12.Container) string {
	values := make([]string, 0, len(containers))
	for _, v := range containers {
		values = append(values, v.Image)
	}
	return strings.Join(values, ",")
}
func GetAppResourceRequirements(containers []v12.Container, resourcetype string) AppResources {
	parsedAppResources := AppResources{}
	var rawMemory int64 = 0
	switch resourcetype {
	case "limits":
		for _, v := range containers {
			parsedAppResources.CPU += v.Resources.Limits.Cpu().AsApproximateFloat64()
			rawMemory += v.Resources.Limits.Memory().Value()
		}
	case "requests":
		for _, v := range containers {
			parsedAppResources.CPU += v.Resources.Requests.Cpu().AsApproximateFloat64()
			rawMemory += v.Resources.Requests.Memory().Value()
		}
	}
	parsedAppResources.Memory = ParseToMB(rawMemory)
	return parsedAppResources
}
func ParseToMB(val int64) string {
	if val > 0 {
		return strconv.FormatInt(val/(1000*1000), 10) + "M"
	}
	return ""
}

func GenericLookupSharedComponents(resourceType string, category string, name string) string {
	if cachedResource, found := GetFromCache("Shared" + resourceType + "Component/" + category + "/" + strings.ToLower(name)); found {
		return cachedResource.(string)
	}
	return ""
}
func GenericUpsertSharedComponents(resourceType string, category string, name string) string {
	if name == "" {
		return ""
	}
	component := ComponentRequest{
		Name:          strings.ToLower(name),
		RootWorkspace: workspaceId,
		TypeID:        lookUpTypeId("Shared" + resourceType + "Component"),
		Fields: map[string]interface{}{
			"shared_category": category,
		},
	}
	componentId := GenericLookupSharedComponents(resourceType, category, name)
	if componentId == "" {
		requestStarted := time.Now()
		resp, err := RestyClient().SetBody(BodyProvider{
			request: component,
			fields:  component.Fields,
		}.Body()).SetResult(&Component{}).Post("component")
		metrics.RequestLatency.WithLabelValues("create").Observe(time.Since(requestStarted).Seconds())
		if err != nil {
			metrics.RequestStatusCode.WithLabelValues("error").Inc()
			log.Errorf("Error creating Shared Components: %s", err)
		}
		cmp := resp.Result().(*Component)
		metrics.RequestStatusCode.WithLabelValues("success").Inc()
		componentId = cmp.ID
		PersistToCache("Shared"+resourceType+"Component/"+category+"/"+strings.ToLower(name), componentId)
		log.Infof("Added Shared Component:%s: %s: %s", resourceType, component.Name, componentId)
		return componentId
	}
	return componentId
}
func GenericDeleteSharedComponents(resourceType string, category string, name string) error {
	var err error
	componentId := GenericLookupSharedComponents(resourceType, category, name)
	if componentId == "" {
		return errors.New("resource not found")
	}
	requestStarted := time.Now()
	_, err = RestyClient().Delete("component/" + componentId)
	metrics.RequestLatency.WithLabelValues("delete").Observe(time.Since(requestStarted).Seconds())
	if err != nil {
		metrics.RequestStatusCode.WithLabelValues("error").Inc()
		log.Errorf("Error deleting Shared%sComponent|%s : %s", resourceType, name, err)
		return err
	}
	metrics.RequestStatusCode.WithLabelValues("success").Inc()
	Cache.Delete("Shared" + resourceType + "Component/" + category + "/" + strings.ToLower(name))
	log.Infof("Deleted Shared%sComponent: %s", resourceType, name)
	return nil
}
func (r *Resource) Link(linkType string, compId string, reverse ...bool) {
	if _, found := GetFromCache("SharedResourceLinks/" + r.ID + "/" + compId); !found && compId != "" {
		referenceLink := ReferenceRequest{
			DisplayText:     linkType,
			RootWorkspace:   workspaceId,
			TargetWorkspace: workspaceId,
			Type:            2,
			Source:          compId,
			Target:          r.ID,
		}
		referenceLink.Description = r.ID + "/" + compId
		if !(len(reverse) > 0 && reverse[0]) {
			referenceLink.Source = r.ID
			referenceLink.Target = compId
		}
		resp, err := RestyClient().SetBody(BodyProvider{
			request: referenceLink,
			fields:  referenceLink.Fields,
		}.Body()).SetResult(&Reference{}).Post("reference")
		if err != nil {
			metrics.RequestStatusCode.WithLabelValues("error").Inc()
			log.Errorf("Error linking resource to a shared component: %s", err)
		}
		reference := resp.Result().(*Reference)
		if reference.ID != "" {
			PersistToCache("SharedResourceLinks/"+r.ID+"/"+compId, reference.ID)
		}
	}

}
func (n *Node) Link(linkType string, compId string, reverse ...bool) {
	if _, found := GetFromCache("SharedNodeLinks/" + n.ID + "/" + compId); !found && compId != "" {
		referenceLink := ReferenceRequest{
			DisplayText:     linkType,
			RootWorkspace:   workspaceId,
			TargetWorkspace: workspaceId,
			Type:            2,
			Source:          compId,
			Target:          n.ID,
		}
		referenceLink.Description = n.ID + "/" + compId
		if !(len(reverse) > 0 && reverse[0]) {
			referenceLink.Source = n.ID
			referenceLink.Target = compId
		}
		resp, err := RestyClient().SetBody(BodyProvider{
			request: referenceLink,
			fields:  referenceLink.Fields,
		}.Body()).SetResult(&Reference{}).Post("reference")
		if err != nil {
			metrics.RequestStatusCode.WithLabelValues("error").Inc()
			log.Errorf("Error linking node to a shared component: %s", err)

		}
		reference := resp.Result().(*Reference)
		if reference.ID != "" {
			PersistToCache("SharedNodeLinks/"+n.ID+"/"+compId, reference.ID)
		}
	}

}
