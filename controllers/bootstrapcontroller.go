package controllers

import (
	"context"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/klog/v2"
)

func BootstrapModel() error {
	yamlFile, err := ioutil.ReadFile("bootstrap_models.yaml")
	if err != nil {
		klog.Errorf("yamlFile.Get err #%v ", err)
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
