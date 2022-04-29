package controllers

import (
	"K8SArdoqBridge/app/lib/metrics"
	"context"
	ardoq "github.com/mories76/ardoq-client-go/pkg"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

func BootstrapModel() error {
	yamlFile, err := ioutil.ReadFile("bootstrap_models.yaml")
	if err != nil {
		log.Errorf("yamlFile.Get err #%v ", err)
		return err
	}
	model := ModelRequest{}
	if err != nil {
		log.Error(err)
		return err
	}
	err = yaml.Unmarshal(yamlFile, &model)

	if err != nil {
		log.Errorf("Unmarshal: %v", err)
		return err
	}
	requestStarted := time.Now()
	workspace, err := ardRestClient().Workspaces().Get(context.TODO(), workspaceId)
	metrics.RequestLatency.WithLabelValues("read").Observe(time.Since(requestStarted).Seconds())
	if err != nil {
		metrics.RequestStatusCode.WithLabelValues("error").Inc()
		log.Errorf("Error getting workspace: %s", err)
		return err
	}
	metrics.RequestStatusCode.WithLabelValues("success").Inc()
	//set componentModel to the componentModel from the found workspace
	componentModel := workspace.ComponentModel
	requestStarted = time.Now()
	currentModel, err := ardRestClient().Models().Read(context.TODO(), componentModel)
	metrics.RequestLatency.WithLabelValues("read").Observe(time.Since(requestStarted).Seconds())
	if err != nil {
		metrics.RequestStatusCode.WithLabelValues("error").Inc()
		log.Errorf("Error getting model: %s", err)
		return err
	}
	metrics.RequestStatusCode.WithLabelValues("success").Inc()

	model.ID = currentModel.ID
	err = UpdateModel(componentModel, model)
	if err != nil {
		log.Errorf("Error updating model: %s", err)
		return err
	}

	return nil
}
func BootstrapFields() error {
	yamlFile, err := ioutil.ReadFile("bootstrap_fields.yaml")
	if err != nil {
		log.Errorf("yamlFile.Get err #%v ", err)
		return err
	}
	var fields []FieldRequest
	if err != nil {
		log.Error(err)
		return err
	}
	err = yaml.Unmarshal(yamlFile, &fields)
	if err != nil {
		log.Errorf("Unmarshal: %v", err)
		return err
	}
	requestStarted := time.Now()
	workspace, err := ardRestClient().Workspaces().Get(context.TODO(), workspaceId)
	metrics.RequestLatency.WithLabelValues("read").Observe(time.Since(requestStarted).Seconds())
	if err != nil {
		metrics.RequestStatusCode.WithLabelValues("error").Inc()
		log.Errorf("Error getting workspace: %s", err)
		return err
	}
	metrics.RequestStatusCode.WithLabelValues("success").Inc()
	//set componentModel to the componentModel from the found workspace
	componentModel := workspace.ComponentModel
	err = CreateFields(componentModel, fields)
	if err != nil {
		log.Errorf("Error updating Fields: %s", err)
		return err
	}
	return nil
}
func InitializeCache() error {
	requestStarted := time.Now()
	components, err := ardRestClient().Components().Search(context.TODO(), &ardoq.ComponentSearchQuery{Workspace: workspaceId})
	metrics.RequestLatency.WithLabelValues("search").Observe(time.Since(requestStarted).Seconds())
	if err != nil {
		metrics.RequestStatusCode.WithLabelValues("error").Inc()
		log.Errorf("Error fetching components: %s", err)
		return err
	}
	metrics.RequestStatusCode.WithLabelValues("success").Inc()
	//get the current cluster
	var clusterComponent ardoq.Component
	var nodeComponents []ardoq.Component
	var namespaceComponents []ardoq.Component
	var resourceComponents []ardoq.Component
	var namespaces []string
	for _, v := range *components {
		if v.Type == "Cluster" && v.Name == os.Getenv("ARDOQ_CLUSTER") {
			clusterComponent = v
			PersistToCache("ResourceType/"+v.Type+"/"+v.Name, v.ID)
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
			PersistToCache("ResourceType/"+v.Type+"/"+v.Name, v.ID)
		}
	}
	//get nodes
	for _, v := range *components {
		if v.Type == "Node" && v.Parent == clusterComponent.ID {
			nodeComponents = append(nodeComponents, v)
			node := Node{
				ID:   v.ID,
				Name: v.Name,
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
				OSImage:           v.Fields["node_os_image"].(string),
				Provider:          v.Fields["node_provider"].(string),
				CreationTimestamp: v.Fields["node_creation_timestamp"].(string),
			}
			PersistToCache("ResourceType/"+v.Type+"/"+v.Name, node)
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
			}
			if i, err := strconv.ParseInt(v.Fields["resource_replicas"].(string), 10, 32); err == nil {
				resource.Replicas = int32(i)
			}
			PersistToCache("ResourceType/"+resource.Namespace+"/"+v.Type+"/"+v.Name, resource)
		}
	}
	//get shared components
	for _, v := range *components {
		if Contains([]string{"SharedResourceComponent", "SharedNodeComponent"}, v.Type) {
			PersistToCache(v.Type+"/"+v.Fields["shared_category"].(string)+"/"+strings.ToLower(v.Name), v.ID)
		}
	}
	requestStarted = time.Now()
	references, err := ardRestClient().References().GetAll(context.TODO())
	metrics.RequestLatency.WithLabelValues("search").Observe(time.Since(requestStarted).Seconds())
	if err != nil {
		metrics.RequestStatusCode.WithLabelValues("error").Inc()
		log.Errorf("Error fetching references: %s", err)
		return err
	}
	metrics.RequestStatusCode.WithLabelValues("success").Inc()
	//get shared references
	for _, v := range *references {
		if Contains(ApplicationLinks, v.DisplayText) && v.RootWorkspace == workspaceId {
			PersistToCache("SharedResourceLinks/"+v.Description, v.ID)
		}
		if Contains(NodeLinks, v.DisplayText) && v.RootWorkspace == workspaceId {
			PersistToCache("SharedNodeLinks/"+v.Description, v.ID)
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
