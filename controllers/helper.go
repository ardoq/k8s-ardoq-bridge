package controllers

import (
	"bytes"
	"context"
	"fmt"
	ardoq "github.com/mories76/ardoq-client-go/pkg"
	goCache "github.com/patrickmn/go-cache"
	"io"
	"io/ioutil"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	"net/http"
	"os"
	"strings"
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
func StripBrackets(in string) string {
	replacer := strings.NewReplacer("[\"", "", "\"]", "")
	return replacer.Replace(in)
}
func GenericLookup(resourceType string, name string, deletion ...bool) string {
	if cachedResource, found := Cache.Get("ResourceType/" + resourceType + "/" + name); found {
		return cachedResource.(string)
	}
	data, err := AdvancedSearch("component", resourceType, name)
	if err != nil {
		klog.Error(err)
		os.Exit(1)
	}
	var componentId string
	if data.Path("total").Data().(float64) == 0 {
		if !(len(deletion) > 0 && deletion[0]) {
			componentId = GenericUpsert(resourceType, name)
		}
		return componentId
	}
	componentId = StripBrackets(data.Search("results", "doc", "_id").String())
	return componentId
}

func lookUpTypeId(name string) string {
	if typeId, found := Cache.Get("ArdoqTypes/" + name); found {
		return typeId.(string)
	}
	workspace, err := ardRestClient().Workspaces().Get(context.TODO(), workspaceId)
	if err != nil {
		klog.Errorf("Error getting workspace: %s", err)
	}
	//set componentModel to the componentModel from the found workspace
	componentModel := workspace.ComponentModel
	model, err := ardRestClient().Models().Read(context.TODO(), componentModel)
	if err != nil {
		klog.Errorf("Error getting model: %s", err)
	}
	cmpTypes := model.GetComponentTypeID()
	if cmpTypes[name] != "" {
		Cache.Set("ArdoqTypes/"+name, cmpTypes[name], goCache.NoExpiration)
		return cmpTypes[name]
	} else {
		return ""
	}

}

func AdvancedSearch(searchType string, resourceType string, name string) (*gabs.Container, error) {
	url := fmt.Sprintf("%sadvanced-search?size=1&from=0", baseUri)
	method := "POST"
	searchQuery := []byte(fmt.Sprintf(`{
			"condition": "AND",
			"rules": [
				{
					"id": "type",
					"field": "type",
					"type": "string",
					"input": "select",
					"operator": "equal",
					"value": "%s"
				},
				{
					"condition": "AND",
					"rules": [
						{
							"id": "rootWorkspace",
							"field": "rootWorkspace",
							"type": "string",
							"input": "text",
							"operator": "equal",
							"value": "%s"
						},
						{
							"id": "typeName",
							"field": "typeName",
							"type": "string",
							"input": "text",
							"operator": "contains",
							"value": "%s"
						},
						{
							"id": "name",
							"field": "name",
							"type": "string",
							"input": "text",
							"operator": "contains",
							"value": "%s"
						}
					]
				}
			]
		}`, searchType, workspaceId, resourceType, name))
	payload := bytes.NewBuffer(searchQuery)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		klog.Fatal(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		klog.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			klog.Fatal(err)
		}
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		klog.Fatal(err)
	}
	parsed, err := gabs.ParseJSON(body)
	if err != nil {
		klog.Fatal(err)
	}
	return parsed, nil
}
func ApplicationResourceSearch(namespace string, resourceType string, resourceName string, deletion ...bool) (*gabs.Container, error) {
	url := fmt.Sprintf("%sadvanced-search?size=1&from=0", baseUri)
	method := "POST"
	parentId := ""
	if len(deletion) > 0 && deletion[0] {
		parentId = GenericLookup("Namespace", namespace, deletion[0])
	} else {
		parentId = GenericLookup("Namespace", namespace)
	}

	searchQuery := []byte(fmt.Sprintf(`{
			"condition": "AND",
			"rules": [
				{
					"id": "type",
					"field": "type",
					"type": "string",
					"input": "select",
					"operator": "equal",
					"value": "component"
				},
				{
					"condition": "AND",
					"rules": [
						{
							"id": "rootWorkspace",
							"field": "rootWorkspace",
							"type": "string",
							"input": "text",
							"operator": "equal",
							"value": "%s"
						},
						{
							"id": "parent",
							"field": "parent",
							"type": "string",
							"input": "text",
							"operator": "equal",
							"value": "%s"
						},
						{
							"id": "typeName",
							"field": "typeName",
							"type": "string",
							"input": "text",
							"operator": "contains",
							"value": "%s"
						},
						{
							"id": "name",
							"field": "name",
							"type": "string",
							"input": "text",
							"operator": "contains",
							"value": "%s"
						}
					]
				}
			]
		}`, workspaceId, parentId, resourceType, resourceName))
	payload := bytes.NewBuffer(searchQuery)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		klog.Fatal(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		klog.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			klog.Fatal(err)
		}
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		klog.Fatal(err)
	}
	parsed, err := gabs.ParseJSON(body)
	if err != nil {
		klog.Fatal(err)
	}
	return parsed, nil
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
			Cache.Set("ResourceType/"+v.Type+"/"+v.Name, v.ID, goCache.NoExpiration)
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
			Cache.Set("ResourceType/"+v.Type+"/"+v.Name, v.ID, goCache.NoExpiration)
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
			Cache.Set("ResourceType/"+v.Type+"/"+v.Name, node, goCache.NoExpiration)
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
			Cache.Set("ResourceType/"+resource.Namespace+"/"+v.Type+"/"+v.Name, resource, goCache.NoExpiration)
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
