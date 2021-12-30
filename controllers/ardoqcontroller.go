package controllers

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Jeffail/gabs"
	ardoq "github.com/mories76/ardoq-client-go/pkg"
	"io"
	"io/ioutil"
	"k8s.io/klog"
	"net/http"
	"os"
	"strings"
)

var (
	baseUri     = os.Getenv("ARDOQ_BASEURI")
	apiKey      = os.Getenv("ARDOQ_APIKEY")
	org         = os.Getenv("ARDOQ_ORG")
	workspaceId = os.Getenv("ARDOQ_WORKSPACE_ID")
	cluster     = os.Getenv("ARDOQ_CLUSTER")
)

type Resource struct {
	RType     string
	Name      string
	ID        string
	Namespace string
	Replicas  int64
	Image     string
}

func ardRestClient() *ardoq.APIClient {
	a, err := ardoq.NewRestClient(baseUri, apiKey, org, "v0.0.0")
	if err != nil {
		fmt.Printf("cannot create new restclient %s", err)
		os.Exit(1)
	}
	return a
}
func stripBrackets(in string) string {
	replacer := strings.NewReplacer("[\"", "", "\"]", "")
	return replacer.Replace(in)
}
func lookUpTypeId(name string) string {
	workspace, err := ardRestClient().Workspaces().Get(context.TODO(), workspaceId)
	if err != nil {
		klog.Error("error getting workspace: %s", err)
	}
	//set componentModel to the componentModel from the found workspace
	componentModel := workspace.ComponentModel
	model, err := ardRestClient().Models().Read(context.TODO(), componentModel)
	if err != nil {
		klog.Error("error getting model: %s", err)
	}
	cmpTypes := model.GetComponentTypeID()
	if cmpTypes[name] != "" {
		return cmpTypes[name]
	} else {
		return ""
	}
}

func advancedSearch(searchType string, queryTypeName string, queryString string) (*gabs.Container, error) {
	url := fmt.Sprintf("%sadvanced-search?size=100&from=0", baseUri)
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
		}`, searchType, workspaceId, queryTypeName, queryString))
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
func LookupCluster(name string) string {
	data, err := advancedSearch("component", "Cluster", name)
	if err != nil {
		klog.Error(err)
		os.Exit(1)
	}
	var componentId string
	if data.Path("total").Data().(float64) == 0 {
		componentId = UpsertCluster(name)
		return componentId
	}
	componentId = stripBrackets(data.Search("results", "doc", "_id").String())
	return componentId
}

func UpsertCluster(name string) string {
	data, err := advancedSearch("component", "Cluster", name)
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
			klog.Errorf("error creating Cluster: %s", err)
		}
		componentId = cmp.ID
		klog.Infof("Cluster: %q: %s", component.Name, componentId)
		return componentId
	}
	componentId = stripBrackets(data.Search("results", "doc", "_id").String())
	_, err = ardRestClient().Components().Update(context.TODO(), componentId, component)
	if err != nil {
		klog.Errorf("error updating Cluster: %s", err)
	}
	klog.Infof("Cluster: %q: %s", component.Name, componentId)
	return componentId
}
func UpsertNamespace(name string) string {
	data, err := advancedSearch("component", "Namespace", name)
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
			klog.Errorf("error creating Namespace: %s", err)
		}
		componentId = cmp.ID
		klog.Infof("Namespace: %q: %s", component.Name, componentId)
		return componentId
	}
	componentId = stripBrackets(data.Search("results", "doc", "_id").String())
	_, err = ardRestClient().Components().Update(context.TODO(), componentId, component)
	if err != nil {
		klog.Errorf("error updating Namespace: %s", err)
	}
	klog.Infof("Namespace: %q: %s", component.Name, componentId)
	return componentId
}
func UpsertDeploymentStatefulset(resource Resource) string {
	data, err := advancedSearch("component", resource.RType, resource.Name)
	if err != nil {
		klog.Error(err)
		os.Exit(1)
	}
	var componentId string
	component := ardoq.ComponentRequest{
		Name:          resource.Name,
		RootWorkspace: workspaceId,
		Parent:        UpsertNamespace(resource.Namespace),
		TypeID:        lookUpTypeId(resource.RType),
		Fields: map[string]interface{}{
			"tags":           resource.RType,
			"resource_image": resource.Image,
			"replicas":       resource.Replicas,
		},
	}
	if data.Path("total").Data().(float64) == 0 {
		cmp, err := ardRestClient().Components().Create(context.TODO(), component)
		if err != nil {
			klog.Errorf("error creating %s : %s", resource.RType, err)
		}
		componentId = cmp.ID
		klog.Infof("%s: %q: %s", resource.RType, resource.Name, componentId)
		return componentId
	}
	componentId = stripBrackets(data.Search("results", "doc", "_id").String())
	_, err = ardRestClient().Components().Update(context.TODO(), componentId, component)
	if err != nil {
		klog.Errorf("error updating %s : %s", resource.RType, err)
	}
	klog.Infof("%s: %q: %s", resource.RType, resource.Name, componentId)
	return componentId
}
func DeleteDeploymentStatefulset(resource Resource) {
	data, err := advancedSearch("component", resource.RType, resource.Name)
	if err != nil {
		klog.Error(err)
		os.Exit(1)
	}
	var componentId string
	if data.Path("total").Data().(float64) == 0 {
		klog.Errorf("error deleting,does not exist: %s", err)
		return
	}
	componentId = stripBrackets(data.Search("results", "doc", "_id").String())
	err = ardRestClient().Components().Delete(context.TODO(), componentId)
	if err != nil {
		klog.Errorf("error deleting %s : %s", resource.RType, err)
	}
	klog.Infof("%s: %q", resource.RType, resource.Name)
	return
}
