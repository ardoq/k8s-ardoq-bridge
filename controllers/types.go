package controllers

import "encoding/json"

type AppResources struct {
	CPU    float64
	Memory string
}
type Resource struct {
	ID                string
	Name              string
	ResourceType      string
	Namespace         string
	Replicas          int32
	Image             string
	CreationTimestamp string
	Stack             string
	Team              string
	Project           string
	Requests          AppResources
	Limits            AppResources
}
type NodeResources struct {
	CPU     int64
	Memory  string
	Storage string
	Pods    int64
}
type Node struct {
	ID                string
	Name              string
	Architecture      string
	Capacity          NodeResources
	Allocatable       NodeResources
	ContainerRuntime  string
	InstanceType      string
	KernelVersion     string
	KubeletVersion    string
	KubeProxyVersion  string
	OperatingSystem   string
	OSImage           string
	Pool              string
	Provider          string
	CreationTimestamp string
	Region            string
	Zone              string
}
type HttpError struct {
	ID, Message string
}
type Model struct {
	Description string                 `json:"description"`
	ID          string                 `json:"_id"`
	Name        string                 `json:"name"`
	Root        ModelComponentTypes    `json:"root"`
	Fields      map[string]interface{} `mapstructure:",remain"`
}
type ModelRequest struct {
	ID          string              `json:"_id"`
	Description string              `json:"description"`
	Root        ModelComponentTypes `json:"root"`
}

type ModelComponentTypes map[string]struct {
	Children     ModelComponentTypes `json:"children"`
	Name         string              `json:"name"`
	ID           string              `json:"id"`
	Index        int                 `json:"index"`
	Icon         string              `json:"icon"`
	Color        string              `json:"color"`
	Image        interface{}         `json:"image"`
	Level        int                 `json:"level"`
	ReturnsValue bool                `json:"returnsValue"`
	Shape        interface{}         `json:"shape"`
	Standard     interface{}         `json:"standard"`
}

type FieldRequest struct {
	ComponentType []string `yaml:"componentType,flow" json:"componentType"`
	Label         string   `json:"label"`
	Model         string   `json:"model"`
	Type          string   `json:"type"`
}
type Component struct {
	Children      []string               `json:"children"`
	Description   string                 `json:"description"`
	ID            string                 `json:"_id" mapstructure:"_id"`
	Model         string                 `json:"model"`
	Name          string                 `json:"name"`
	Parent        string                 `json:"parent,omitempty"`
	RootWorkspace string                 `json:"rootWorkspace"`
	Type          string                 `json:"type"`
	TypeID        string                 `json:"typeId" mapstructure:"typeId"`
	Fields        map[string]interface{} `json:"-" mapstructure:",remain"`
}

// MarshalJSON implements custom JSON marshaling for Component
// Fields are flattened to the top level to match the Ardoq API format
func (c Component) MarshalJSON() ([]byte, error) {
	// Create a map with standard fields
	result := make(map[string]interface{})

	if len(c.Children) > 0 {
		result["children"] = c.Children
	}
	if c.Description != "" {
		result["description"] = c.Description
	}
	if c.ID != "" {
		result["_id"] = c.ID
	}
	if c.Model != "" {
		result["model"] = c.Model
	}
	if c.Name != "" {
		result["name"] = c.Name
	}
	if c.Parent != "" {
		result["parent"] = c.Parent
	}
	if c.RootWorkspace != "" {
		result["rootWorkspace"] = c.RootWorkspace
	}
	if c.Type != "" {
		result["type"] = c.Type
	}
	if c.TypeID != "" {
		result["typeId"] = c.TypeID
	}

	// Flatten custom fields to the top level
	if c.Fields != nil {
		for key, value := range c.Fields {
			// Skip standard fields to avoid duplication
			if key != "children" && key != "description" && key != "_id" &&
				key != "model" && key != "name" && key != "parent" &&
				key != "rootWorkspace" && key != "type" && key != "typeId" {
				result[key] = value
			}
		}
	}

	return json.Marshal(result)
}

type ComponentRequest struct {
	RootWorkspace string                 `json:"rootWorkspace,omitempty"`
	Name          interface{}            `json:"name,omitempty"`
	Description   interface{}            `json:"description,omitempty"`
	Parent        interface{}            `json:"parent,omitempty"`
	TypeID        interface{}            `json:"typeId,omitempty"`
	Fields        map[string]interface{} `json:"-"`
}
type Reference struct {
	DisplayText     string                 `json:"displayText"`
	Description     string                 `json:"description"`
	ID              string                 `json:"_id"`
	RootWorkspace   string                 `json:"rootWorkspace"`
	Source          string                 `json:"source"`
	Target          string                 `json:"target"`
	TargetWorkspace string                 `json:"targetWorkspace"`
	Type            int                    `json:"type"`
	Model           string                 `json:"model"`
	Version         int                    `json:"_version"`
	Fields          map[string]interface{} `mapstructure:",remain"`
}

type ReferenceRequest struct {
	Description     interface{}            `json:"description,omitempty"`
	DisplayText     interface{}            `json:"displayText,omitempty"`
	RootWorkspace   interface{}            `json:"rootWorkspace,omitempty"`
	Source          interface{}            `json:"source,omitempty"`
	Target          interface{}            `json:"target,omitempty"`
	TargetWorkspace interface{}            `json:"targetWorkspace,omitempty"`
	Type            interface{}            `json:"type,omitempty"`
	Fields          map[string]interface{} `json:"-"`
}
type Workspace struct {
	ComponentModel string                 `json:"componentModel"`
	Fields         map[string]interface{} `json:",remain"`
	ID             string                 `json:"_id"`
	Name           string                 `json:"name"`
}

// GetComponentTypeID recursively traverses the model component types and returns a map of component type names to their IDs
func (m *Model) GetComponentTypeID() map[string]string {
	result := make(map[string]string)
	m.traverseComponentTypes(m.Root, result)
	return result
}

// traverseComponentTypes is a helper function that recursively walks through ModelComponentTypes
func (m *Model) traverseComponentTypes(types ModelComponentTypes, result map[string]string) {
	for _, componentType := range types {
		if componentType.Name != "" && componentType.ID != "" {
			result[componentType.Name] = componentType.ID
		}
		if componentType.Children != nil {
			m.traverseComponentTypes(componentType.Children, result)
		}
	}
}
