package kubecontroller

import (
	brucov1alpha1 "github.com/ferama/bruco/pkg/kube/apis/brucocontroller/v1alpha1"
	"gopkg.in/yaml.v2"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func newConfigmap(bruco *brucov1alpha1.Bruco) *corev1.ConfigMap {
	b, _ := yaml.Marshal(bruco.Spec.Conf)
	confMap := make(map[string]string)

	confMap["config.yaml"] = string(b)
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      bruco.Name,
			Namespace: bruco.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(bruco, brucov1alpha1.SchemeGroupVersion.WithKind("Bruco")),
			},
		},
		Data: confMap,
	}
}

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
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "bruco-config",
									MountPath: "/bruco",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "bruco-config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: bruco.Name,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
