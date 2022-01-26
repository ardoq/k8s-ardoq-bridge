package subscriptions

import (
	"K8SArdoqBridge/app/controllers"
	"context"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getNamespaceLabels(name string) map[string]string {
	namespace, err := controllers.ClientSet.CoreV1().Namespaces().Get(context.TODO(), name, v12.GetOptions{})
	if err != nil {
		return nil
	}
	return namespace.GetLabels()
}
