package controllers

type Resource struct {
	ID           string
	Name         string
	ResourceType string
	Namespace    string
	Replicas     int64
	Image        string
}
type NodeResources struct {
	CPU     int64
	Memory  string
	Storage string
	Pods    int64
}
type Node struct {
	ID               string
	Name             string
	Architecture     string
	Capacity         NodeResources
	Allocatable      NodeResources
	ContainerRuntime string
	KernelVersion    string
	KubeletVersion   string
	KubeProxyVersion string
	OperatingSystem  string
	OSImage          string
	Provider         string
}