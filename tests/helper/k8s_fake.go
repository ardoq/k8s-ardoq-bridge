package helper

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

// NewFakeK8sClient creates a new fake Kubernetes clientset for testing
func NewFakeK8sClient() kubernetes.Interface {
	return fake.NewSimpleClientset()
}

// DeploymentOptions configures a test Deployment
type DeploymentOptions struct {
	Name      string
	Namespace string
	Replicas  int32
	Image     string
	Labels    map[string]string
	PodLabels map[string]string
}

// CreateFakeDeployment creates a Deployment in the fake clientset
func CreateFakeDeployment(client kubernetes.Interface, opts DeploymentOptions) (*appsv1.Deployment, error) {
	if opts.Namespace == "" {
		opts.Namespace = "default"
	}
	if opts.Image == "" {
		opts.Image = "nginx:1.14.2"
	}
	if opts.Replicas == 0 {
		opts.Replicas = 2
	}
	if opts.Labels == nil {
		opts.Labels = make(map[string]string)
	}
	if opts.PodLabels == nil {
		opts.PodLabels = map[string]string{"app": "nginx"}
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      opts.Name,
			Namespace: opts.Namespace,
			Labels:    opts.Labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &opts.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: opts.PodLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: opts.PodLabels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: opts.Image,
							Ports: []corev1.ContainerPort{
								{ContainerPort: 80},
							},
						},
					},
				},
			},
		},
	}

	return client.AppsV1().Deployments(opts.Namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
}

// StatefulSetOptions configures a test StatefulSet
type StatefulSetOptions struct {
	Name      string
	Namespace string
	Replicas  int32
	Image     string
	Labels    map[string]string
	PodLabels map[string]string
}

// CreateFakeStatefulSet creates a StatefulSet in the fake clientset
func CreateFakeStatefulSet(client kubernetes.Interface, opts StatefulSetOptions) (*appsv1.StatefulSet, error) {
	if opts.Namespace == "" {
		opts.Namespace = "default"
	}
	if opts.Image == "" {
		opts.Image = "nginx:1.14.2"
	}
	if opts.Replicas == 0 {
		opts.Replicas = 2
	}
	if opts.Labels == nil {
		opts.Labels = make(map[string]string)
	}
	if opts.PodLabels == nil {
		opts.PodLabels = map[string]string{"app": "nginx"}
	}

	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      opts.Name,
			Namespace: opts.Namespace,
			Labels:    opts.Labels,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas:    &opts.Replicas,
			ServiceName: "nginx",
			Selector: &metav1.LabelSelector{
				MatchLabels: opts.PodLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: opts.PodLabels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: opts.Image,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
									Name:          "web",
								},
							},
						},
					},
				},
			},
		},
	}

	return client.AppsV1().StatefulSets(opts.Namespace).Create(context.TODO(), statefulSet, metav1.CreateOptions{})
}

// NamespaceOptions configures a test Namespace
type NamespaceOptions struct {
	Name   string
	Labels map[string]string
}

// CreateFakeNamespace creates a Namespace in the fake clientset
func CreateFakeNamespace(client kubernetes.Interface, opts NamespaceOptions) (*corev1.Namespace, error) {
	if opts.Labels == nil {
		opts.Labels = make(map[string]string)
	}

	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   opts.Name,
			Labels: opts.Labels,
		},
	}

	return client.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
}

// NodeOptions configures a test Node
type NodeOptions struct {
	Name        string
	Labels      map[string]string
	CPUCapacity string
	MemCapacity string
}

// CreateFakeNode creates a Node in the fake clientset
func CreateFakeNode(client kubernetes.Interface, opts NodeOptions) (*corev1.Node, error) {
	if opts.Labels == nil {
		opts.Labels = make(map[string]string)
	}
	if opts.CPUCapacity == "" {
		opts.CPUCapacity = "4"
	}
	if opts.MemCapacity == "" {
		opts.MemCapacity = "8Gi"
	}

	node := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:   opts.Name,
			Labels: opts.Labels,
		},
		Status: corev1.NodeStatus{
			Capacity: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse(opts.CPUCapacity),
				corev1.ResourceMemory: resource.MustParse(opts.MemCapacity),
			},
			Allocatable: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse(opts.CPUCapacity),
				corev1.ResourceMemory: resource.MustParse(opts.MemCapacity),
			},
			NodeInfo: corev1.NodeSystemInfo{
				Architecture:            "amd64",
				ContainerRuntimeVersion: "docker://20.10.0",
				KernelVersion:           "5.10.0",
				KubeletVersion:          "v1.24.0",
				KubeProxyVersion:        "v1.24.0",
				OperatingSystem:         "linux",
				OSImage:                 "Ubuntu 20.04",
			},
		},
	}

	return client.CoreV1().Nodes().Create(context.TODO(), node, metav1.CreateOptions{})
}

// DeleteFakeDeployment deletes a Deployment from the fake clientset
func DeleteFakeDeployment(client kubernetes.Interface, namespace, name string) error {
	return client.AppsV1().Deployments(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

// DeleteFakeStatefulSet deletes a StatefulSet from the fake clientset
func DeleteFakeStatefulSet(client kubernetes.Interface, namespace, name string) error {
	return client.AppsV1().StatefulSets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

// DeleteFakeNamespace deletes a Namespace from the fake clientset
func DeleteFakeNamespace(client kubernetes.Interface, name string) error {
	return client.CoreV1().Namespaces().Delete(context.TODO(), name, metav1.DeleteOptions{})
}

// DeleteFakeNode deletes a Node from the fake clientset
func DeleteFakeNode(client kubernetes.Interface, name string) error {
	return client.CoreV1().Nodes().Delete(context.TODO(), name, metav1.DeleteOptions{})
}
