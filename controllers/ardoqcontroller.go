package controllers

import (
	"context"
	"errors"
	"github.com/Jeffail/gabs"
	ardoq "github.com/mories76/ardoq-client-go/pkg"
	goCache "github.com/patrickmn/go-cache"
	"k8s.io/klog/v2"
	"os"
	"reflect"
)

var (
	baseUri                       = os.Getenv("ARDOQ_BASEURI")
	apiKey                        = os.Getenv("ARDOQ_APIKEY")
	org                           = os.Getenv("ARDOQ_ORG")
	workspaceId                   = os.Getenv("ARDOQ_WORKSPACE_ID")
	cluster                       = os.Getenv("ARDOQ_CLUSTER")
	validApplicationResourceTypes = []string{"Deployment", "StatefulSet"}
)

func GenericUpsert(resourceType string, genericResource interface{}) string {
	var (
		data     *gabs.Container
		err      error
		resource Resource
		node     Node
		name     string
	)
	if Contains([]string{"Cluster", "Namespace"}, resourceType) {
		name = genericResource.(string)
		data, err = AdvancedSearch("component", resourceType, name)
	} else if Contains(validApplicationResourceTypes, resourceType) {
		resource = genericResource.(Resource)
		name = resource.Name
		data, err = ApplicationResourceSearch(resource.Namespace, resource.ResourceType, name)
	} else if resourceType == "Node" {
		node = genericResource.(Node)
		name = node.Name
		data, err = AdvancedSearch("component", resourceType, name)
	} else {
		err = errors.New("invalid resource type")
	}

	if err != nil {
		klog.Error(err)
		os.Exit(1)
	}
	var componentId string
	component := ardoq.ComponentRequest{
		Name:          name,
		RootWorkspace: workspaceId,
		TypeID:        lookUpTypeId(resourceType),
	}
	switch resourceType {
	case "Namespace":
		component.Parent = GenericLookup("Cluster", cluster)
		break
	case "Deployment", "StatefulSet":
		component.Parent = GenericLookup("Namespace", resource.Namespace)
		component.Fields = map[string]interface{}{
			"tags":           resource.ResourceType,
			"resource_image": resource.Image,
			"replicas":       resource.Replicas,
		}
		break
	case "Node":
		component.Parent = GenericLookup("Cluster", cluster)
		component.Fields = map[string]interface{}{
			"node_architecture":        node.Architecture,
			"node_container_runtime":   node.ContainerRuntime,
			"node_kernel_version":      node.KernelVersion,
			"node_kubelet_version":     node.KubeletVersion,
			"node_kube_proxy_version":  node.KubeProxyVersion,
			"node_os":                  node.OperatingSystem,
			"node_os_image":            node.OSImage,
			"node_capacity_cpu":        node.Capacity.CPU,
			"node_capacity_memory":     node.Capacity.Memory,
			"node_capacity_storage":    node.Capacity.Storage,
			"node_capacity_pods":       node.Capacity.Pods,
			"node_allocatable_cpu":     node.Allocatable.CPU,
			"node_allocatable_memory":  node.Allocatable.Memory,
			"node_allocatable_storage": node.Allocatable.Storage,
			"node_allocatable_pods":    node.Allocatable.Pods,
			"node_provider":            node.Provider,
		}
		break
	}
	if data.Path("total").Data().(float64) == 0 {
		cmp, err := ardRestClient().Components().Create(context.TODO(), component)
		if err != nil {
			klog.Errorf("Error creating %s: %s", resourceType, err)
		}
		componentId = cmp.ID
		switch resourceType {
		case "Namespace", "Cluster":
			Cache.Set("ResourceType/"+resourceType+"/"+name, componentId, goCache.NoExpiration)
			break
		case "Deployment", "StatefulSet":
			resource.ID = componentId
			Cache.Set("ResourceType/"+resource.Namespace+"/"+resourceType+"/"+name, resource, goCache.NoExpiration)
			break
		case "Node":
			node.ID = componentId
			Cache.Set("ResourceType/"+resourceType+"/"+name, node, goCache.NoExpiration)
			break
		}
		klog.Infof("Added %s: %q: %s", resourceType, component.Name, componentId)
		ApplyDelay()
		return componentId
	}
	componentId = StripBrackets(data.Search("results", "doc", "_id").String())
	switch resourceType {
	case "Namespace", "Cluster":
		if cachedResource, found := Cache.Get("ResourceType/" + resourceType + "/" + name); found {
			return cachedResource.(string)
		} else {
			Cache.Set("ResourceType/"+resourceType+"/"+name, componentId, goCache.NoExpiration)
		}
		break
	case "Deployment", "StatefulSet":
		resource.ID = componentId
		if cachedResource, found := Cache.Get("ResourceType/" + resource.Namespace + "/" + resourceType + "/" + name); found && cachedResource.(Resource) == resource {
			return componentId
		}
		Cache.Set("ResourceType/"+resource.Namespace+"/"+resourceType+"/"+name, resource, goCache.NoExpiration)
		break
	case "Node":
		node.ID = componentId
		if cachedResource, found := Cache.Get("ResourceType/" + resourceType + "/" + name); found && cachedResource.(Node) == node {
			return componentId
		}
		Cache.Set("ResourceType/"+resourceType+"/"+name, node, goCache.NoExpiration)
		break
	}
	_, err = ardRestClient().Components().Update(context.TODO(), componentId, component)
	if err != nil {
		klog.Errorf("Error updating %s: %s", resourceType, err)
	}
	klog.Infof("Updated %s: %q: %s", resourceType, component.Name, componentId)
	return componentId
}
func GenericDelete(resourceType string, genericResource interface{}) error {
	var (
		data     *gabs.Container
		err      error
		resource Resource
		node     Node
		name     string
	)
	if Contains([]string{"Cluster", "Namespace"}, resourceType) {
		name = reflect.ValueOf(genericResource).String()
		data, err = AdvancedSearch("component", resourceType, name)
	} else if Contains(validApplicationResourceTypes, resourceType) {
		resource = genericResource.(Resource)
		name = resource.Name
		data, err = ApplicationResourceSearch(resource.Namespace, resource.ResourceType, name)
	} else if resourceType == "Node" {
		node = genericResource.(Node)
		name = node.Name
		data, err = AdvancedSearch("component", resourceType, name)
	} else {
		err = errors.New("invalid resource type")
	}

	if err != nil {
		klog.Error(err)
	}
	var componentId string
	if data.Path("total").Data().(float64) == 0 {
		return errors.New("resource not found")
	}
	componentId = StripBrackets(data.Search("results", "doc", "_id").String())
	err = ardRestClient().Components().Delete(context.TODO(), componentId)
	if err != nil {
		klog.Errorf("Error deleting %s : %s", resourceType, err)
		return err
	}
	switch resourceType {
	case "Deployment", "StatefulSet":
		Cache.Delete("ResourceType/" + resource.Namespace + "/" + resourceType + "/" + name)
		break
	default:
		Cache.Delete("ResourceType/" + resourceType + "/" + name)
		break
	}

	klog.Infof("Deleted %s: %q", resourceType, name)
	return nil
}
