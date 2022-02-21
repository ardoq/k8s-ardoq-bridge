package controllers

import (
	"K8SArdoqBridge/app/lib/metrics"
	"context"
	"errors"
	ardoq "github.com/mories76/ardoq-client-go/pkg"
	"k8s.io/klog/v2"
	"os"
	"time"
)

var (
	baseUri                       = os.Getenv("ARDOQ_BASEURI")
	apiKey                        = os.Getenv("ARDOQ_APIKEY")
	org                           = os.Getenv("ARDOQ_ORG")
	workspaceId                   = os.Getenv("ARDOQ_WORKSPACE_ID")
	validApplicationResourceTypes = []string{"Deployment", "StatefulSet"}
)

func GenericUpsert(resourceType string, genericResource interface{}) string {
	var (
		componentId string
		err         error
		resource    Resource
		node        Node
		name        string
	)
	switch resourceType {
	case "Namespace", "Cluster":
		name = genericResource.(string)
		break
	case "Deployment", "StatefulSet":
		resource = genericResource.(Resource)
		name = resource.Name
		break
	case "Node":
		node = genericResource.(Node)
		name = node.Name
		break
	default:
		err = errors.New("invalid resource type")
	}

	if err != nil {
		klog.Error(err)
		os.Exit(1)
	}
	component := ardoq.ComponentRequest{
		Name:          name,
		RootWorkspace: workspaceId,
		TypeID:        lookUpTypeId(resourceType),
	}
	switch resourceType {
	case "Namespace":
		component.Parent = LookupCluster(os.Getenv("ARDOQ_CLUSTER"))
		componentId = LookupNamespace(name)
		break
	case "Deployment", "StatefulSet":
		namespace := LookupNamespace(resource.Namespace)
		if namespace == "" {
			namespace = GenericUpsert("Namespace", resource.Namespace)
		}
		component.Parent = namespace
		component.Fields = map[string]interface{}{
			"resource_image":              resource.Image,
			"resource_replicas":           resource.Replicas,
			"resource_creation_timestamp": resource.CreationTimestamp,
			"resource_stack":              resource.Stack,
			"resource_team":               resource.Team,
			"resource_project":            resource.Project,
		}
		componentId = LookupResource(resource.Namespace, resourceType, name)
		break
	case "Node":
		component.Parent = LookupCluster(os.Getenv("ARDOQ_CLUSTER"))
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
			"node_creation_timestamp":  node.CreationTimestamp,
			"node_zone":                node.Zone,
			"node_region":              node.Region,
		}
		componentId = LookupNode(name)
		break
	}
	if componentId == "" {
		requestStarted := time.Now()
		cmp, err := ardRestClient().Components().Create(context.TODO(), component)
		metrics.RequestLatency.WithLabelValues("create").Observe(time.Since(requestStarted).Seconds())
		if err != nil {
			klog.Errorf("Error creating %s: %s", resourceType, err)
		}
		componentId = cmp.ID
		switch resourceType {
		case "Namespace", "Cluster":
			PersistToCache("ResourceType/"+resourceType+"/"+name, componentId)
			break
		case "Deployment", "StatefulSet":
			resource.ID = componentId
			PersistToCache("ResourceType/"+resource.Namespace+"/"+resourceType+"/"+name, resource)
			break
		case "Node":
			node.ID = componentId
			PersistToCache("ResourceType/"+resourceType+"/"+name, node)
			break
		}
		klog.Infof("Added %s: %q: %s", resourceType, component.Name, componentId)
		return componentId
	}
	switch resourceType {
	case "Namespace", "Cluster":
		if cachedResource, found := GetFromCache("ResourceType/" + resourceType + "/" + name); found {
			return cachedResource.(string)
		}
		PersistToCache("ResourceType/"+resourceType+"/"+name, componentId)
		break
	case "Deployment", "StatefulSet":
		resource.ID = componentId
		if cachedResource, found := GetFromCache("ResourceType/" + resource.Namespace + "/" + resourceType + "/" + name); found && cachedResource.(Resource) == resource {
			return componentId
		}
		PersistToCache("ResourceType/"+resource.Namespace+"/"+resourceType+"/"+name, resource)
		break
	case "Node":
		node.ID = componentId
		if cachedResource, found := GetFromCache("ResourceType/" + resourceType + "/" + name); found && cachedResource.(Node) == node {
			return componentId
		}
		PersistToCache("ResourceType/"+resourceType+"/"+name, node)
		break
	}
	requestStarted := time.Now()
	_, err = ardRestClient().Components().Update(context.TODO(), componentId, component)
	metrics.RequestLatency.WithLabelValues("update").Observe(time.Since(requestStarted).Seconds())
	if err != nil {
		klog.Errorf("Error updating %s|%s: %s", resourceType, name, err)
	}
	klog.Infof("Updated %s: %q: %s", resourceType, component.Name, componentId)
	return componentId
}
func GenericDelete(resourceType string, genericResource interface{}) error {
	var (
		componentId string
		err         error
		resource    Resource
		node        Node
		name        string
	)
	switch resourceType {
	case "Cluster":
		name = genericResource.(string)
		componentId = LookupCluster(name, true)
		break
	case "Namespace":
		name = genericResource.(string)
		componentId = LookupNamespace(name)
		break
	case "Deployment", "StatefulSet":
		resource = genericResource.(Resource)
		name = resource.Name
		componentId = LookupResource(resource.Namespace, resourceType, name)
		break
	case "Node":
		node = genericResource.(Node)
		name = node.Name
		componentId = LookupNode(name)
		break
	default:
		err = errors.New("invalid resource type")
	}

	if err != nil {
		klog.Error(err)
	}
	if componentId == "" {
		return errors.New("resource not found")
	}
	requestStarted := time.Now()
	err = ardRestClient().Components().Delete(context.TODO(), componentId)
	metrics.RequestLatency.WithLabelValues("delete").Observe(time.Since(requestStarted).Seconds())
	if err != nil {
		klog.Errorf("Error deleting %s|%s : %s", resourceType, name, err)
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
