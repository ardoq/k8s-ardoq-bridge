package controllers

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
