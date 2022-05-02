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
	ID                  string                 `json:"_id"`
	LastUpdated2        string                 `json:"last-updated"`
	LastModifiedBy      string                 `json:"last-modified-by"`
	LastModifiedByEmail string                 `json:"lastModifiedByEmail"`
	LastModifiedByName  string                 `json:"lastModifiedByName"`
	LastUpdated         string                 `json:"lastupdated"`
	Model               string                 `json:"model"`
	Name                string                 `json:"name"`
	Order               float64                `json:"_order"`
	Parent              interface{}            `json:"parent"`
	RootWorkspace       string                 `json:"rootWorkspace"`
	Type                string                 `json:"type"`
	TypeID              string                 `json:"typeId"`
	Version             int                    `json:"_version"`
	Fields              map[string]interface{} `json:",remain"`
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
	Fields              map[string]interface{} `json:",remain"`
}

// ReferenceRequest is the payload for creating and updating a reference
// Fields map has json tag "-" so that it doesn't get marshalled into JSON
// the fields are being handled by the ardoqBodyProvider
// URL: PATCH/POST /api/reference
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
