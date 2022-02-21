package controllers

import (
	"context"
	ardoq "github.com/mories76/ardoq-client-go/pkg"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/klog/v2"
	"os"
	"strconv"
)

func BootstrapModel() error {
	yamlFile, err := ioutil.ReadFile("bootstrap_models.yaml")
	if err != nil {
		klog.Errorf("yamlFile.Get err #%v ", err)
		return err
	}
	model := ModelRequest{}
	if err != nil {
		klog.Error(err)
		return err
	}
	err = yaml.Unmarshal(yamlFile, &model)

	if err != nil {
		klog.Errorf("Unmarshal: %v", err)
		return err
	}

	workspace, err := ardRestClient().Workspaces().Get(context.TODO(), workspaceId)
	if err != nil {
		klog.Errorf("Error getting workspace: %s", err)
		return err
	}
	//set componentModel to the componentModel from the found workspace
	componentModel := workspace.ComponentModel
	currentModel, err := ardRestClient().Models().Read(context.TODO(), componentModel)
	if err != nil {
		klog.Errorf("Error getting model: %s", err)
		return err
	}

	model.ID = currentModel.ID
	err = UpdateModel(context.TODO(), componentModel, model)
	if err != nil {
		klog.Errorf("Error updating model: %s", err)
		return err
	}

	return nil
}
func BootstrapFields() error {
	yamlFile, err := ioutil.ReadFile("bootstrap_fields.yaml")
	if err != nil {
		klog.Errorf("yamlFile.Get err #%v ", err)
		return err
	}
	var fields []FieldRequest
	if err != nil {
		klog.Error(err)
		return err
	}
	err = yaml.Unmarshal(yamlFile, &fields)
	if err != nil {
		klog.Errorf("Unmarshal: %v", err)
		return err
	}
	workspace, err := ardRestClient().Workspaces().Get(context.TODO(), workspaceId)
	if err != nil {
		klog.Errorf("Error getting workspace: %s", err)
		return err
	}
	//set componentModel to the componentModel from the found workspace
	componentModel := workspace.ComponentModel
	err = CreateFields(context.TODO(), componentModel, fields)
	if err != nil {
		klog.Errorf("Error updating Fields: %s", err)
		return err
	}
	return nil
}
func InitializeCache() error {
	components, err := ardRestClient().Components().Search(context.TODO(), &ardoq.ComponentSearchQuery{Workspace: workspaceId})
	if err != nil {
		klog.Errorf("Error fetching components %s: %s", err)
		return err
	}
	//get the current cluster
	var clusterComponent ardoq.Component
	var nodeComponents []ardoq.Component
	var namespaceComponents []ardoq.Component
	var resourceComponents []ardoq.Component
	var namespaces []string
	for _, v := range *components {
		if v.Type == "Cluster" && v.Name == os.Getenv("ARDOQ_CLUSTER") {
			clusterComponent = v
			SetToCache("ResourceType/"+v.Type+"/"+v.Name, v.ID)
		}
	}
	if clusterComponent.ID == "" {
		return nil
	}
	//get namespaces
	for _, v := range *components {
		if v.Type == "Namespace" && v.Parent == clusterComponent.ID {
			namespaceComponents = append(namespaceComponents, v)
			namespaces = append(namespaces, v.ID)
			SetToCache("ResourceType/"+v.Type+"/"+v.Name, v.ID)
		}
	}
	//get nodes
	for _, v := range *components {
		if v.Type == "Node" && v.Parent == clusterComponent.ID {
			nodeComponents = append(nodeComponents, v)
			node := Node{
				ID:           v.ID,
				Name:         v.Name,
				Architecture: v.Fields["node_architecture"].(string),
				Capacity: NodeResources{
					CPU:     int64(v.Fields["node_capacity_cpu"].(float64)),
					Memory:  v.Fields["node_capacity_memory"].(string),
					Storage: v.Fields["node_capacity_storage"].(string),
					Pods:    int64(v.Fields["node_allocatable_pods"].(float64)),
				},
				Allocatable: NodeResources{
					CPU:     int64(v.Fields["node_allocatable_cpu"].(float64)),
					Memory:  v.Fields["node_allocatable_memory"].(string),
					Storage: v.Fields["node_allocatable_storage"].(string),
					Pods:    int64(v.Fields["node_allocatable_pods"].(float64)),
				},
				ContainerRuntime:  v.Fields["node_container_runtime"].(string),
				KernelVersion:     v.Fields["node_kernel_version"].(string),
				KubeletVersion:    v.Fields["node_kubelet_version"].(string),
				KubeProxyVersion:  v.Fields["node_kube_proxy_version"].(string),
				OperatingSystem:   v.Fields["node_os"].(string),
				OSImage:           v.Fields["node_os_image"].(string),
				Provider:          v.Fields["node_provider"].(string),
				CreationTimestamp: v.Fields["node_creation_timestamp"].(string),
				Region:            v.Fields["node_zone"].(string),
				Zone:              v.Fields["node_region"].(string),
			}
			SetToCache("ResourceType/"+v.Type+"/"+v.Name, node)
		}
	}
	//get application resources
	for _, v := range *components {
		if Contains([]string{"Deployment", "StatefulSet"}, v.Type) && Contains(namespaces, v.Parent.(string)) {
			resourceComponents = append(resourceComponents, v)
			resource := Resource{
				ID:                v.ID,
				Name:              v.Name,
				ResourceType:      v.Type,
				Namespace:         getNamespace(namespaceComponents, v.Parent.(string)),
				Image:             v.Fields["resource_image"].(string),
				CreationTimestamp: v.Fields["resource_creation_timestamp"].(string),
				Stack:             v.Fields["resource_stack"].(string),
				Team:              v.Fields["resource_team"].(string),
				Project:           v.Fields["resource_project"].(string),
			}
			if i, err := strconv.ParseInt(v.Fields["resource_replicas"].(string), 10, 32); err == nil {
				resource.Replicas = int32(i)
			}
			SetToCache("ResourceType/"+resource.Namespace+"/"+v.Type+"/"+v.Name, resource)
		}
	}
	return nil
}
func getNamespace(namespaceComponents []ardoq.Component, id string) string {
	for _, v := range namespaceComponents {
		if v.ID == id {
			return v.Name
		}
	}
	return ""

}
