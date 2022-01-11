package controllers

import (
	"context"
	"errors"
	ardoq "github.com/mories76/ardoq-client-go/pkg"
	"k8s.io/klog"
	"os"
)

var (
	baseUri            = os.Getenv("ARDOQ_BASEURI")
	apiKey             = os.Getenv("ARDOQ_APIKEY")
	org                = os.Getenv("ARDOQ_ORG")
	workspaceId        = os.Getenv("ARDOQ_WORKSPACE_ID")
	cluster            = os.Getenv("ARDOQ_CLUSTER")
	validResourceTypes = []string{"Deployment", "StatefulSet"}
)

func UpsertCluster(name string) string {
	data, err := AdvancedSearch("component", "Cluster", name)
	if err != nil {
		klog.Error(err)
		os.Exit(1)
	}
	var componentId string
	component := ardoq.ComponentRequest{
		Name:          name,
		RootWorkspace: workspaceId,
		TypeID:        lookUpTypeId("Cluster"),
		Fields: map[string]interface{}{
			"Tags": "Cluster",
		},
	}
	if data.Path("total").Data().(float64) == 0 {
		cmp, err := ardRestClient().Components().Create(context.TODO(), component)
		if err != nil {
			klog.Errorf("Error creating Cluster: %s", err)
		}
		componentId = cmp.ID
		klog.Infof("Added Cluster: %q: %s", component.Name, componentId)
		return componentId
	}
	componentId = StripBrackets(data.Search("results", "doc", "_id").String())
	_, err = ardRestClient().Components().Update(context.TODO(), componentId, component)
	if err != nil {
		klog.Errorf("Error updating Cluster: %s", err)
	}
	klog.Infof("Updated Cluster: %q: %s", component.Name, componentId)
	return componentId
}
func DeleteCluster(cluster string) error {
	data, err := AdvancedSearch("component", "Cluster", cluster)
	if err != nil {
		klog.Error(err)
		os.Exit(1)
	}
	var componentId string
	if data.Path("total").Data().(float64) == 0 {
		return errors.New("cluster not found")
	}
	componentId = StripBrackets(data.Search("results", "doc", "_id").String())
	err = ardRestClient().Components().Delete(context.TODO(), componentId)
	if err != nil {
		klog.Errorf("Error deleting Cluster : %s", err)
		return err
	}
	klog.Infof("Deleted Cluster: %q", cluster)
	return nil
}
func UpsertNamespace(name string) string {
	data, err := AdvancedSearch("component", "Namespace", name)
	if err != nil {
		klog.Error(err)
		os.Exit(1)
	}
	var componentId string
	component := ardoq.ComponentRequest{
		Name:          name,
		RootWorkspace: workspaceId,
		Parent:        LookupCluster(cluster),
		TypeID:        lookUpTypeId("Namespace"),
		Fields: map[string]interface{}{
			"tags": "Namespace",
		},
	}
	if data.Path("total").Data().(float64) == 0 {
		cmp, err := ardRestClient().Components().Create(context.TODO(), component)
		if err != nil {
			klog.Errorf("Error creating Namespace: %s", err)
		}
		componentId = cmp.ID
		klog.Infof("Added Namespace: %q: %s", component.Name, componentId)
		ApplyDelay()
		return componentId
	}
	componentId = StripBrackets(data.Search("results", "doc", "_id").String())
	_, err = ardRestClient().Components().Update(context.TODO(), componentId, component)
	if err != nil {
		klog.Errorf("Error updating Namespace: %s", err)
	}
	klog.Infof("Updated Namespace: %q: %s", component.Name, componentId)
	return componentId
}
func DeleteNamespace(cluster string) error {
	data, err := AdvancedSearch("component", "Namespace", cluster)
	if err != nil {
		klog.Error(err)
		os.Exit(1)
	}
	var componentId string
	if data.Path("total").Data().(float64) == 0 {
		return errors.New("namespace not found")
	}
	componentId = StripBrackets(data.Search("results", "doc", "_id").String())
	err = ardRestClient().Components().Delete(context.TODO(), componentId)
	if err != nil {
		klog.Errorf("Error deleting Namespace : %s", err)
		return err
	}
	klog.Infof("Deleted Namespace: %q", cluster)
	return nil
}
func UpsertApplicationResource(resource Resource) string {
	//data, err := AdvancedSearch("component", resource.ResourceType, resource.Name)
	data, err := ApplicationResourceSearch(resource.Namespace, resource.ResourceType, resource.Name)
	if err != nil {
		klog.Error(err)
		os.Exit(1)
	}
	var componentId string
	component := ardoq.ComponentRequest{
		Name:          resource.Name,
		RootWorkspace: workspaceId,
		Parent:        UpsertNamespace(resource.Namespace),
		TypeID:        lookUpTypeId(resource.ResourceType),
		Fields: map[string]interface{}{
			"tags":           resource.ResourceType,
			"resource_image": resource.Image,
			"replicas":       resource.Replicas,
		},
	}
	if data.Path("total").Data().(float64) == 0 {
		cmp, err := ardRestClient().Components().Create(context.TODO(), component)
		if err != nil {
			klog.Errorf("Error creating %s : %s", resource.ResourceType, err)
		}
		componentId = cmp.ID
		klog.Infof("Added %s: %q: %s", resource.ResourceType, resource.Name, componentId)
		return componentId
	}
	componentId = StripBrackets(data.Search("results", "doc", "_id").String())
	_, err = ardRestClient().Components().Update(context.TODO(), componentId, component)
	if err != nil {
		klog.Errorf("Error updating %s : %s", resource.ResourceType, err)
	}
	klog.Infof("Updated %s: %q: %s", resource.ResourceType, resource.Name, componentId)
	return componentId
}
func DeleteApplicationResource(resource Resource) error {
	//data, err := AdvancedSearch("component", resource.ResourceType, resource.Name)
	data, err := ApplicationResourceSearch(resource.Namespace, resource.ResourceType, resource.Name)
	if err != nil {
		klog.Error(err)
		os.Exit(1)
	}
	var componentId string
	if data.Path("total").Data().(float64) == 0 {
		return errors.New("resource not found")
	}
	componentId = StripBrackets(data.Search("results", "doc", "_id").String())
	err = ardRestClient().Components().Delete(context.TODO(), componentId)
	if err != nil {
		klog.Errorf("Error deleting %s : %s", resource.ResourceType, err)
		return err
	}
	klog.Infof("Deleted %s: %q", resource.ResourceType, resource.Name)
	return nil
}
func UpsertNode(node Node) string {
	data, err := AdvancedSearch("component", "Node", node.Name)
	if err != nil {
		klog.Error(err)
		os.Exit(1)
	}
	var componentId string
	component := ardoq.ComponentRequest{
		Name:          node.Name,
		RootWorkspace: workspaceId,
		Parent:        LookupCluster(cluster),
		TypeID:        lookUpTypeId("Node"),
		Fields: map[string]interface{}{
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
		},
	}
	if data.Path("total").Data().(float64) == 0 {
		cmp, err := ardRestClient().Components().Create(context.TODO(), component)
		if err != nil {
			klog.Errorf("Error creating Node: %s", err)
		}
		componentId = cmp.ID
		klog.Infof("Added Node: %q: %s", component.Name, componentId)
		return componentId
	}
	componentId = StripBrackets(data.Search("results", "doc", "_id").String())
	_, err = ardRestClient().Components().Update(context.TODO(), componentId, component)
	if err != nil {
		klog.Errorf("Error updating Node: %s", err)
	}
	klog.Infof("Updated Node: %q: %s", component.Name, componentId)
	return componentId
}
func DeleteNode(node Node) error {
	data, err := AdvancedSearch("component", "Node", node.Name)
	if err != nil {
		klog.Error(err)
		os.Exit(1)
	}
	var componentId string
	if data.Path("total").Data().(float64) == 0 {
		return errors.New("node not found")
	}
	componentId = StripBrackets(data.Search("results", "doc", "_id").String())
	err = ardRestClient().Components().Delete(context.TODO(), componentId)
	if err != nil {
		klog.Errorf("Error deleting Node : %s", err)
		return err
	}
	klog.Infof("Deleted Node: %q", node.Name)
	return nil
}
