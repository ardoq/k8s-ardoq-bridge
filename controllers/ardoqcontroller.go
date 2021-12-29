package controllers

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Jeffail/gabs"
	ardoq "github.com/mories76/ardoq-client-go/pkg"
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

type searchResponse gabs.Container

type Resoure struct {
	Name string
	ID   string
}

func ardRestClient() *ardoq.APIClient {
	a, err := ardoq.NewRestClient(baseUri, apiKey, org, "v0.0.0")
	if err != nil {
		//return nil, errors.Wrap(err, "cannot create new restclient")
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
		klog.Error("error during get workspace: %s", err)
	}
	//set componentModel to the componentModel from the found workspace
	componentModel := workspace.ComponentModel
	model, err := ardRestClient().Models().Read(context.TODO(), componentModel)
	if err != nil {
		klog.Error("error during get model: %s", err)
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
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		klog.Fatal(err)
	}
	//var searchRes searchResponse
	//err2 := json.Unmarshal(body, &searchRes)
	parsed, err2 := gabs.ParseJSON(body)
	if err2 != nil {
		klog.Fatal(err2)
	}
	return parsed, nil
}
func firstOrCreateCluster(name string) string {
	data, err := advancedSearch("component", "Cluster", name)
	var componentId string
	if err != nil {
		klog.Error(err)
		os.Exit(1)
	}
	if data.Path("total").Data().(float64) == 0 {
		component := ardoq.ComponentRequest{
			Name:          name,
			RootWorkspace: workspaceId,
			TypeID:        lookUpTypeId("Cluster"),
		}
		cmp, err := ardRestClient().Components().Create(context.TODO(), component)
		if err != nil {
			klog.Error("error during component create: %s", err)
		}
		componentId = cmp.ID
		klog.Info("Cluster: " + componentId)
		return componentId
	}
	componentId = stripBrackets(data.Search("results", "doc", "_id").String())
	klog.Info("Cluster: " + componentId)
	return componentId
}
func firstOrCreateNamespace(name string) string {
	data, err := advancedSearch("component", "Namespace", name)
	var componentId string
	if err != nil {
		klog.Error(err)
		os.Exit(1)
	}
	if data.Path("total").Data().(float64) == 0 {
		component := ardoq.ComponentRequest{
			Name:          name,
			RootWorkspace: workspaceId,
			Parent:        firstOrCreateCluster(cluster),
			TypeID:        lookUpTypeId("Namespace"),
		}
		cmp, err := ardRestClient().Components().Create(context.TODO(), component)
		if err != nil {
			klog.Error("error during component create: %s", err)
		}
		componentId = cmp.ID
		klog.Info("Namespace: " + componentId)
		return componentId
	}
	componentId = stripBrackets(data.Search("results", "doc", "_id").String())
	klog.Info("Namespace: " + componentId)
	return componentId
}
func firstOrCreateDeploymentStatefulset(name string, resouceType string, namespace string) string {
	data, err := advancedSearch("component", resouceType, name)
	var componentId string
	if err != nil {
		klog.Error(err)
		os.Exit(1)
	}
	if data.Path("total").Data().(float64) == 0 {
		component := ardoq.ComponentRequest{
			Name:          name,
			RootWorkspace: workspaceId,
			Parent:        firstOrCreateNamespace(namespace),
			TypeID:        lookUpTypeId(resouceType),
		}
		cmp, err := ardRestClient().Components().Create(context.TODO(), component)
		if err != nil {
			klog.Error("error during component create: %s", err)
		}
		componentId = cmp.ID
		klog.Info(resouceType + ": " + componentId)
		return componentId
	}
	componentId = stripBrackets(data.Search("results", "doc", "_id").String())
	klog.Info(resouceType + ": " + componentId)
	return componentId
}
