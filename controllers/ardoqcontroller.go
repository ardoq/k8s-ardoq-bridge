package controllers

import (
	"K8SArdoqBridge/app/lib/metrics"
	"errors"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	validApplicationResourceTypes = []string{"Deployment", "StatefulSet"}
	ApplicationLinks              = []string{"Owns", "Is realized By", "Is Supported By"}
	NodeLinks                     = []string{"LOCATION", "SUB_LOCATION", "ARCHITECTURE", "INSTANCE_TYPE", "OS", "NODE_POOL"}
)

// Configuration getters - read from environment on each call to support testing
func getBaseUri() string     { return os.Getenv("ARDOQ_BASEURI") }
func getApiKey() string      { return os.Getenv("ARDOQ_APIKEY") }
func getOrg() string         { return os.Getenv("ARDOQ_ORG") }
func getWorkspaceId() string { return os.Getenv("ARDOQ_WORKSPACE_ID") }

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
		log.Error(err)
		os.Exit(1)
	}
	component := ComponentRequest{
		Name:          name,
		RootWorkspace: getWorkspaceId(),
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
			"resource_requests_cpu":       resource.Requests.CPU,
			"resource_requests_memory":    resource.Requests.Memory,
			"resource_limits_cpu":         resource.Limits.CPU,
			"resource_limits_memory":      resource.Limits.Memory,
		}
		componentId = LookupResource(resource.Namespace, resourceType, name)
		break
	case "Node":
		component.Parent = LookupCluster(os.Getenv("ARDOQ_CLUSTER"))
		component.Fields = map[string]interface{}{
			"node_container_runtime":   node.ContainerRuntime,
			"node_kernel_version":      node.KernelVersion,
			"node_kubelet_version":     node.KubeletVersion,
			"node_kube_proxy_version":  node.KubeProxyVersion,
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
		}
		componentId = LookupNode(name)
		break
	}
	if componentId == "" {
		requestStarted := time.Now()
		resp, err := RestyClient().SetBody(BodyProvider{
			request: component,
			fields:  component.Fields,
		}.Body()).SetResult(&Component{}).Post("component")
		metrics.RequestLatency.WithLabelValues("create").Observe(time.Since(requestStarted).Seconds())
		if err != nil {
			metrics.RequestStatusCode.WithLabelValues("error").Inc()
			log.Errorf("Error creating %s: %s", resourceType, err)
		}
		cmp := resp.Result().(*Component)
		metrics.RequestStatusCode.WithLabelValues("success").Inc()
		componentId = cmp.ID
		switch resourceType {
		case "Namespace", "Cluster":
			PersistToCache("ResourceType/"+resourceType+"/"+name, componentId)
			break
		case "Deployment", "StatefulSet":
			resource.ID = componentId
			PersistToCache("ResourceType/"+resource.Namespace+"/"+resourceType+"/"+name, resource)
			resource.Link("Owns", GenericUpsertSharedComponents("Resource", "team", resource.Team), true)
			resource.Link("Is realized By", GenericUpsertSharedComponents("Resource", "stack", resource.Stack))
			resource.Link("Is Supported By", GenericUpsertSharedComponents("Resource", "project", resource.Project))
			break
		case "Node":
			node.ID = componentId
			PersistToCache("ResourceType/"+resourceType+"/"+name, node)
			node.Link("LOCATION", GenericUpsertSharedComponents("Node", "location", node.Region))
			node.Link("SUB_LOCATION", GenericUpsertSharedComponents("Node", "sub_location", node.Zone))
			node.Link("ARCHITECTURE", GenericUpsertSharedComponents("Node", "architecture", node.Architecture))
			node.Link("INSTANCE_TYPE", GenericUpsertSharedComponents("Node", "instance_type", node.InstanceType))
			node.Link("OS", GenericUpsertSharedComponents("Node", "node_os", node.OperatingSystem))
			node.Link("NODE_POOL", GenericUpsertSharedComponents("Node", "node_pool", node.Pool))
			break
		}
		log.Infof("Added %s: %s: %s", resourceType, component.Name, componentId)
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
		resource.Link("Owns", GenericUpsertSharedComponents("Resource", "team", resource.Team), true)
		resource.Link("Is realized By", GenericUpsertSharedComponents("Resource", "stack", resource.Stack))
		resource.Link("Is Supported By", GenericUpsertSharedComponents("Resource", "project", resource.Project))
		break
	case "Node":
		node.ID = componentId
		if cachedResource, found := GetFromCache("ResourceType/" + resourceType + "/" + name); found && cachedResource.(Node) == node {
			return componentId
		}
		PersistToCache("ResourceType/"+resourceType+"/"+name, node)
		node.Link("LOCATION", GenericUpsertSharedComponents("Node", "location", node.Region))
		node.Link("SUB_LOCATION", GenericUpsertSharedComponents("Node", "sub_location", node.Zone))
		node.Link("ARCHITECTURE", GenericUpsertSharedComponents("Node", "architecture", node.Architecture))
		node.Link("INSTANCE_TYPE", GenericUpsertSharedComponents("Node", "instance_type", node.InstanceType))
		node.Link("OS", GenericUpsertSharedComponents("Node", "node_os", node.OperatingSystem))
		node.Link("NODE_POOL", GenericUpsertSharedComponents("Node", "node_pool", node.Pool))
		break
	}
	requestStarted := time.Now()
	_, err = RestyClient().SetBody(BodyProvider{
		request: component,
		fields:  component.Fields,
	}.Body()).SetResult(&Component{}).Patch("component/" + componentId)
	metrics.RequestLatency.WithLabelValues("update").Observe(time.Since(requestStarted).Seconds())
	if err != nil {
		metrics.RequestStatusCode.WithLabelValues("error").Inc()
		log.Errorf("Error updating %s|%s: %s", resourceType, name, err)
	}
	metrics.RequestStatusCode.WithLabelValues("success").Inc()
	log.Infof("Updated %s: %s: %s", resourceType, component.Name, componentId)
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
		log.Error(err)
	}
	if componentId == "" {
		return errors.New("resource not found")
	}
	requestStarted := time.Now()
	_, err = RestyClient().Delete("component/" + componentId)
	metrics.RequestLatency.WithLabelValues("delete").Observe(time.Since(requestStarted).Seconds())
	if err != nil {
		metrics.RequestStatusCode.WithLabelValues("error").Inc()
		log.Errorf("Error deleting %s|%s : %s", resourceType, name, err)
		return err
	}
	metrics.RequestStatusCode.WithLabelValues("success").Inc()
	switch resourceType {
	case "Deployment", "StatefulSet":
		Cache.Delete("ResourceType/" + resource.Namespace + "/" + resourceType + "/" + name)
		break
	default:
		Cache.Delete("ResourceType/" + resourceType + "/" + name)
		break
	}
	log.Infof("Deleted %s: %s", resourceType, name)
	return nil
}
