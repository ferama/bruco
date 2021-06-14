package kubecontroller

import (
	brucov1alpha1 "github.com/ferama/bruco/pkg/kube/apis/brucocontroller/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func newService(bruco *brucov1alpha1.Bruco) *corev1.Service {
	labels := map[string]string{
		"app":        "bruco",
		"controller": bruco.Name,
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      bruco.Name,
			Namespace: bruco.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(bruco, brucov1alpha1.SchemeGroupVersion.WithKind("Bruco")),
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port: 8080,
					TargetPort: intstr.IntOrString{
						IntVal: 8080,
					},
				},
			},
			Selector: labels,
		},
	}
}

// newDeployment creates a new Deployment for a Bruco resource. It also sets
// the appropriate OwnerReferences on the resource so handleObject can discover
// the Bruco resource that 'owns' it.
func newDeployment(bruco *brucov1alpha1.Bruco) *appsv1.Deployment {
	labels := map[string]string{
		"app":        "bruco",
		"controller": bruco.Name,
	}
	containerImage := "ferama/bruco:dev"
	if bruco.Spec.ContainerImage != "" {
		containerImage = bruco.Spec.ContainerImage
	}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      bruco.Name,
			Namespace: bruco.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(bruco, brucov1alpha1.SchemeGroupVersion.WithKind("Bruco")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: bruco.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "bruco",
							Image:   containerImage,
							Command: []string{"bruco", bruco.Spec.FunctionURL},
							Env:     bruco.Spec.Env,
						},
					},
				},
			},
		},
	}
}
