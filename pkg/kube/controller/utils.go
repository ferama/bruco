package kubecontroller

import (
	"fmt"

	brucov1alpha1 "github.com/ferama/bruco/pkg/kube/apis/brucocontroller/v1alpha1"
)

func getDeploymentName(bruco *brucov1alpha1.Bruco) string {
	depName := bruco.Name
	if bruco.Spec.Name != "" {
		depName = bruco.Spec.Name
	}
	return fmt.Sprintf("bruco-%s", depName)
}

func getServiceName(bruco *brucov1alpha1.Bruco) string {
	return fmt.Sprintf("bruco-%s", bruco.Name)
}

func getConfigMapName(bruco *brucov1alpha1.Bruco) string {
	return fmt.Sprintf("bruco-%s", bruco.Name)
}
