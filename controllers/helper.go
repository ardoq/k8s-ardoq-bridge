package controllers

import (
	"K8SArdoqBridge/app/lib/metrics"
	"context"
	"fmt"
	ardoq "github.com/mories76/ardoq-client-go/pkg"
	goCache "github.com/patrickmn/go-cache"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	"os"
	"time"
)

var (
	Cache     = goCache.New(5*time.Minute, 10*time.Minute)
	ClientSet *kubernetes.Clientset
)

func ardRestClient() *ardoq.APIClient {
	a, err := ardoq.NewRestClient(baseUri, apiKey, org, "v0.0.0")
	if err != nil {
		fmt.Printf("cannot create new restclient %s", err)
		os.Exit(1)
	}
	return a
}
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
		klog.Errorf("Error getting workspace: %s", err)
	}
	//set componentModel to the componentModel from the found workspace
	componentModel := workspace.ComponentModel
	requestStarted = time.Now()
	model, err := ardRestClient().Models().Read(context.TODO(), componentModel)
	metrics.RequestLatency.WithLabelValues("read").Observe(time.Since(requestStarted).Seconds())
	if err != nil {
		klog.Errorf("Error getting model: %s", err)
	}
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
