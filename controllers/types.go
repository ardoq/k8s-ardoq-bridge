package controllers

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
	Ardoq struct {
		EntityType             string      `json:"entity-type"`
		IncomingReferenceCount int         `json:"incomingReferenceCount"`
		OutgoingReferenceCount int         `json:"outgoingReferenceCount"`
		Persistent             interface{} `json:"persistent"`
	} `json:"ardoq"`
	ArdoqPersistent     []interface{} `json:"ardoq-persistent"`
	BlankTemplate       bool          `json:"blankTemplate"`
	Category            string        `json:"category"`
	Common              bool          `json:"common"`
	Created             string        `json:"created"`
	CreatedBy           string        `json:"createdBy"`
	CreatedBy2          string        `json:"created-by"`
	CreatedByEmail      string        `json:"createdByEmail"`
	CreatedByName       string        `json:"createdByName"`
	CreatedFromTemplate string        `json:"createdFromTemplate"`
	DefaultSort         string        `json:"defaultSort"`
	DefaultViews        []string      `json:"defaultViews"`
	Description         string        `json:"description"`
	Flexible            bool          `json:"flexible"`
	Folder              string        `json:"folder"`
	ID                  string        `json:"_id"`
	LastModifiedBy      string        `json:"lastModifiedBy"`
	LastModifiedBy2     string        `json:"last-modified-by"`
	LastModifiedByEmail string        `json:"lastModifiedByEmail"`
	LastModifiedByName  string        `json:"lastModifiedByName"`
	LastUpdated         string        `json:"lastupdated"`
	LastUpdated2        string        `json:"last-updated"`
	MaxReferenceTypeKey int           `json:"maxReferenceTypeKey"`
	Name                string        `json:"name"`
	Origin              struct {
		ID      string `json:"id"`
		Version int    `json:"_version"`
	}
	ReferenceTypes ModelReferenceTypes `json:"referenceTypes"`
	Root           ModelComponentTypes `json:"root"`
	StartView      string              `json:"startView"`
	UseAsTemplate  bool                `json:"useAsTemplate"`
	Version        int                 `json:"_version"`
	Workspaces     struct {
		Restricted int `json:"restricted"`
		UsedBy     []struct {
			ID                  string `json:"_id"`
			Name                string `json:"name"`
			CreatedByName       string `json:"createdByName"`
			CreatedByEmail      string `json:"createdByEmail"`
			LastModifiedByName  string `json:"lastModifiedByName"`
			LastModifiedByEmail string `json:"lastModifiedByEmail"`
			Ardoq               struct {
				EntityType string `json:"entity-type"`
			} `json:"ardoq"`
		} `json:"used-by"`
	} `json:"workspaces"`
	// Fields, is the safetynet for when mapping the API response to the struct.
	// The goal is to have all fields documented in the API to have a known field in the struct.
	// For models, Fields should be null
	Fields map[string]interface{} `json:",remain"`
}
type ModelRequest struct {
	ID          string              `json:"_id"`
	Description string              `json:"description"`
	Root        ModelComponentTypes `json:"root"`
}
type ModelReferenceTypes map[string]struct {
	Name         string `json:"name"`
	ID           int    `json:"id"`
	Color        string `json:"color"`
	Line         string `json:"line"`
	LineEnding   string `json:"lineEnding"`
	ReturnsValue bool   `json:"returnsValue"`
	SvgStyle     string `json:"svgStyle"`
}
type ModelComponentTypes map[string]struct {
	Children     ModelComponentTypes `json:"children"`
	Name         string              `json:"name"`
	ID           string              `json:"id"`
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
	Global        bool     `json:"global"`
	Label         string   `json:"label"`
	Model         string   `json:"model"`
	Type          string   `json:"type"`
}
type Component struct {
	Ardoq struct {
		EntityType             string      `json:"entity-type"`
		Persistent             interface{} `json:"persistent"`
		IncomingReferenceCount int         `json:"incomingReferenceCount"`
		OutgoingReferenceCount int         `json:"outgoingReferenceCount"`
	} `json:"ardoq"`
	Children            []string               `json:"children"`
	ComponentKey        string                 `json:"component-key"`
	Created             string                 `json:"created"`
	CreatedBy           string                 `json:"created-by"`
	CreatedByEmail      string                 `json:"createdByEmail"`
	CreatedByName       string                 `json:"createdByName"`
	Description         string                 `json:"description"`
	ID                  string                 `json:"_id" mapstructure:"_id"`
	LastUpdated2        string                 `json:"last-updated"`
	LastModifiedBy      string                 `json:"last-modified-by"`
	LastModifiedByEmail string                 `json:"lastModifiedByEmail"`
	LastModifiedByName  string                 `json:"lastModifiedByName"`
	LastUpdated         string                 `json:"lastupdated"`
	Model               string                 `json:"model"`
	Name                string                 `json:"name"`
	Order               float64                `json:"_order"`
	Parent              string                 `json:"parent,omitempty"`
	RootWorkspace       string                 `json:"rootWorkspace"`
	Type                string                 `json:"type"`
	TypeID              string                 `json:"typeId"`
	Version             int                    `json:"_version"`
	Fields              map[string]interface{} `mapstructure:",remain"`
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
	Ardoq struct {
		EntityType             string      `json:"entity-type"`
		IncomingReferenceCount int         `json:"incomingReferenceCount"`
		OutgoingReferenceCount int         `json:"outgoingReferenceCount"`
		Persistent             interface{} `json:"persistent"`
	} `json:"ardoq"`
	Created             string                 `json:"created"`
	CreatedBy           string                 `json:"created-by"`
	CreatedByEmail      string                 `json:"createdByEmail"`
	CreatedByName       string                 `json:"createdByName"`
	DisplayText         string                 `json:"displayText"`
	Description         string                 `json:"description"`
	ID                  string                 `json:"_id"`
	LastUpdated2        string                 `json:"last-updated"`
	LastModifiedBy      string                 `json:"last-modified-by"`
	LastModifiedByName  string                 `json:"lastModifiedByName"`
	LastModifiedByEmail string                 `json:"lastModifiedByEmail"`
	LastUpdated         string                 `json:"lastupdated"`
	Order               int                    `json:"order"`
	RootWorkspace       string                 `json:"rootWorkspace"`
	Source              string                 `json:"source"`
	Target              string                 `json:"target"`
	TargetWorkspace     string                 `json:"targetWorkspace"`
	Type                int                    `json:"type"`
	Model               string                 `json:"model"`
	Version             int                    `json:"_version"`
	Fields              map[string]interface{} `mapstructure:",remain"`
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
	Ardoq struct {
		EntityType string `json:"entity-type"`
	} `json:"ardoq"`
	ArdoqPersistent     []interface{}          `json:"ardoq-persistent"`
	CompCounter         int                    `json:"comp-counter"`
	ComponentModel      string                 `json:"componentModel"`
	ComponentTemplate   string                 `json:"componentTemplate"`
	Created             string                 `json:"created"`
	CreatedBy           string                 `json:"created-by"`
	CreatedByEmail      string                 `json:"createdByEmail"`
	CreatedByName       string                 `json:"createdByName"`
	Description         string                 `json:"description"`
	DefaultPerspective  string                 `json:"defaultPerspective"`
	Fields              map[string]interface{} `json:",remain"`
	ID                  string                 `json:"_id"`
	LastUpdated2        string                 `json:"last-updated"`
	LastModifiedBy      string                 `json:"last-modified-by"`
	LastModifiedByEmail string                 `json:"lastModifiedByEmail"`
	LastModifiedByName  string                 `json:"lastModifiedByName"`
	LastUpdated         string                 `json:"lastupdated"`
	LinkedWorkspaces    struct {
		Linked     []string `json:"linked"`
		BackLinked []string `json:"backlinked"`
	} `json:"linked-workspaces"`
	Name   string `json:"name"`
	Origin struct {
		EntityType string `json:"entity-type"`
	} `json:"origin"`
	StartView    string   `json:"startView"`
	Type         string   `json:"type"`
	Version      int      `json:"_version"`
	Views        []string `json:"views"`
	WorkspaceKey string   `json:"workspace-key"`
}
