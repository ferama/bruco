package kubecontroller

import (
	"fmt"

	brucov1alpha1 "github.com/ferama/bruco/pkg/kube/apis/brucocontroller/v1alpha1"
	"gopkg.in/yaml.v2"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	ContainerDefaultImage = "ferama/bruco:dev"
)

// Creates a new config map with bruco config
func newConfigMap(bruco *brucov1alpha1.Bruco) *corev1.ConfigMap {
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

func newFunctionContainer(bruco *brucov1alpha1.Bruco, configName string) *corev1.Container {
	containerImage := ContainerDefaultImage

	if bruco.Spec.Image != "" {
		containerImage = bruco.Spec.Image
	}
	container := &corev1.Container{
		Name:    "bruco",
		Image:   containerImage,
		Command: []string{"bruco", bruco.Spec.FunctionURL},
		Env:     bruco.Spec.Env,
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      configName,
				MountPath: "/bruco",
			},
		},
	}
	if bruco.Spec.ImagePullPolicy == "" {
		container.ImagePullPolicy = corev1.PullAlways
	} else {
		container.ImagePullPolicy = bruco.Spec.ImagePullPolicy
	}

	container.Resources = bruco.Spec.Resources
	return container
}

// newDeployment creates a new Deployment for a Bruco resource. It also sets
// the appropriate OwnerReferences on the resource so handleObject can discover
// the Bruco resource that 'owns' it.
func newDeployment(bruco *brucov1alpha1.Bruco) *appsv1.Deployment {
	labels := map[string]string{
		"app":        "bruco",
		"controller": bruco.Name,
	}
	imagePullSecrets := []corev1.LocalObjectReference{}
	if bruco.Spec.ImagePullSecrets != "" {
		imagePullSecrets = []corev1.LocalObjectReference{
			{Name: bruco.Spec.ImagePullSecrets},
		}
	}
	configName := fmt.Sprintf("bruco-config-%d", bruco.Generation)
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
					ImagePullSecrets: imagePullSecrets,
					Containers: []corev1.Container{
						*newFunctionContainer(bruco, configName),
					},
					Volumes: []corev1.Volume{
						{
							Name: configName,
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
