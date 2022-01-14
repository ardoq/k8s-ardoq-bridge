package controllers

import (
	"K8SArdoqBridge/app/tests/helper"
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
	"time"
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
func LookupCluster(name string) string {
	data, err := AdvancedSearch("component", "Cluster", name)
	if err != nil {
		klog.Error(err)
		os.Exit(1)
	}
	var componentId string
	if data.Path("total").Data().(float64) == 0 {
		componentId = UpsertCluster(name)
		return componentId
	}
	componentId = StripBrackets(data.Search("results", "doc", "_id").String())
	return componentId
}
func LookupNamespace(name string) string {
	data, err := AdvancedSearch("component", "Namespace", name)
	if err != nil {
		klog.Error(err)
		os.Exit(1)
	}
	var componentId string
	if data.Path("total").Data().(float64) == 0 {
		return componentId
	}
	componentId = StripBrackets(data.Search("results", "doc", "_id").String())
	return componentId
}
func lookUpTypeId(name string) string {
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
		return cmpTypes[name]
	} else {
		return ""
	}
}

func AdvancedSearch(searchType string, queryTypeName string, queryString string) (*gabs.Container, error) {
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
func ApplicationResourceSearch(namespace string, resourceType string, resourceName string) (*gabs.Container, error) {
	url := fmt.Sprintf("%sadvanced-search?size=1&from=0", baseUri)
	method := "POST"
	parentId := LookupNamespace(namespace)
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
	if r.Name != "" && r.Namespace != "" && r.ResourceType != "" && r.Image != "" && helper.Contains(validResourceTypes, r.ResourceType) {
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
func ApplyDelay() {
	time.Sleep(2 * time.Second)
}
