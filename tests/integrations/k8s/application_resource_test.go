package k8s_test

import (
	"K8SArdoqBridge/app/controllers"
	"K8SArdoqBridge/app/tests/helper"
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v12 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
	"k8s.io/utils/pointer"
)

var _ = Describe("ApplicationResource", Ordered, func() {
	kubectlClient := helper.KubectlClient()
	var deployResourceName = helper.RandomString(10)
	var stsResourceName = helper.RandomString(10)

	BeforeAll(func() {
		deploy := &v12.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      deployResourceName,
				Namespace: v1.NamespaceDefault,
				Labels: map[string]string{
					"sync-to-ardoq": "true",
				},
			},
			Spec: v12.DeploymentSpec{
				Replicas: pointer.Int32Ptr(2),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": "test",
					},
				},
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app": "test",
						},
					},
					Spec: v1.PodSpec{
						Containers: []v1.Container{
							{
								Name:  "web",
								Image: "nginx:1.14.2",
								Ports: []v1.ContainerPort{
									{
										Name:          "http",
										Protocol:      v1.ProtocolTCP,
										ContainerPort: 80,
									},
								},
							},
						},
					},
				},
			},
		}
		sts := &v12.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      stsResourceName,
				Namespace: v1.NamespaceDefault,
				Labels: map[string]string{
					"sync-to-ardoq": "true",
				},
			},
			Spec: v12.StatefulSetSpec{
				Replicas: pointer.Int32Ptr(2),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": "nginx",
					},
				},
				ServiceName: "nginx",
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app": "nginx",
						},
					},
					Spec: v1.PodSpec{
						Containers: []v1.Container{
							{
								Name:  "nginx",
								Image: "k8s.gcr.io/nginx-slim:0.8",
								Ports: []v1.ContainerPort{
									{
										Name:          "http",
										Protocol:      v1.ProtocolTCP,
										ContainerPort: 80,
									},
								},
								VolumeMounts: []v1.VolumeMount{
									{
										Name:      "www",
										MountPath: "/usr/share/nginx/html",
									},
								},
							},
						},
					},
				},
				VolumeClaimTemplates: []v1.PersistentVolumeClaim{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "www",
						},
						Spec: v1.PersistentVolumeClaimSpec{
							AccessModes: []v1.PersistentVolumeAccessMode{
								"ReadWriteOnce",
							},
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									v1.ResourceStorage: resource.MustParse("1Gi"),
								},
							},
						},
					},
				},
			},
		}
		klog.Info("Creating deployment...")
		deployClient := kubectlClient.AppsV1().Deployments(v1.NamespaceDefault)
		deployResult, err := deployClient.Create(context.TODO(), deploy, metav1.CreateOptions{})
		if err != nil {
			panic(err)
		}
		klog.Infof("Created deployment %q.", deployResult.GetObjectMeta().GetName())
		klog.Info("Creating statefulset...")
		stsClient := kubectlClient.AppsV1().StatefulSets(v1.NamespaceDefault)
		stsResult, err := stsClient.Create(context.TODO(), sts, metav1.CreateOptions{})
		if err != nil {
			panic(err)
		}
		klog.Infof("Created statefulset %q.", stsResult.GetObjectMeta().GetName())
		controllers.ApplyDelay()
	})
	AfterAll(func() {
		deletePolicy := metav1.DeletePropagationForeground
		deployClient := kubectlClient.AppsV1().Deployments(v1.NamespaceDefault)
		if err := deployClient.Delete(context.TODO(), deployResourceName, metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		}); err != nil {
			klog.Error(err)
		}
		klog.Info("Cleaned up deployment.")
		stsClient := kubectlClient.AppsV1().StatefulSets(v1.NamespaceDefault)
		if err := stsClient.Delete(context.TODO(), stsResourceName, metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		}); err != nil {
			klog.Error(err)
		}
		klog.Info("Cleaned up statefulsets.")
	})
	Context("Deployment tests", Ordered, func() {
		It("Can fetch created Deployments", func() {
			data, err := controllers.ApplicationResourceSearch("default", "Deployment", deployResourceName)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data.Path("total").Data().(float64)).ShouldNot(BeZero())
		})
		It("Can not find deleted Deployments", func() {
			klog.Info("Deleting deployment...")
			deletePolicy := metav1.DeletePropagationForeground
			resourceClient := kubectlClient.AppsV1().Deployments(v1.NamespaceDefault)
			if err := resourceClient.Delete(context.TODO(), deployResourceName, metav1.DeleteOptions{
				PropagationPolicy: &deletePolicy,
			}); err != nil {
				panic(err)
			}
			klog.Info("Deleted deployment.")
			controllers.ApplyDelay(10)

			data, err := controllers.ApplicationResourceSearch(v1.NamespaceDefault, "Deployment", deployResourceName)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data.Path("total").Data().(float64)).Should(BeZero())
		})

	})
	Context("StatefulSet tests", Ordered, func() {
		It("Can fetch created StatefulSet", func() {
			data, err := controllers.ApplicationResourceSearch(v1.NamespaceDefault, "StatefulSet", stsResourceName)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data.Path("total").Data().(float64)).ShouldNot(BeZero())
		})
		It("Can not find deleted StatefulSets", func() {
			klog.Info("Deleting statefulset...")
			deletePolicy := metav1.DeletePropagationForeground
			resourceClient := kubectlClient.AppsV1().StatefulSets(v1.NamespaceDefault)
			if err := resourceClient.Delete(context.TODO(), stsResourceName, metav1.DeleteOptions{
				PropagationPolicy: &deletePolicy,
			}); err != nil {
				panic(err)
			}
			klog.Info("Deleted statefulset.")
			controllers.ApplyDelay(5)
			data, err := controllers.ApplicationResourceSearch(v1.NamespaceDefault, "StatefulSet", stsResourceName)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data.Path("total").Data().(float64)).Should(BeZero())
		})
	})
})
