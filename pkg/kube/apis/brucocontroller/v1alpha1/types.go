package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Bruco is a specification for a Bruco resource
type Bruco struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BrucoSpec   `json:"spec"`
	Status BrucoStatus `json:"status"`
}

type BrucoConf map[string]interface{}

func (c *BrucoConf) DeepCopy() *BrucoConf {
	return &BrucoConf{}
}

// BrucoSpec is the spec for a Bruco resource
type BrucoSpec struct {
	Replicas *int32 `json:"replicas"`
	// you may want to use a custom image that has dependencies already installed
	// for example
	Image            string                      `json:"image,omitempty"`
	FunctionURL      string                      `json:"functionURL"`
	ImagePullPolicy  corev1.PullPolicy           `json:"imagePullPolicy,omitempty"`
	ImagePullSecrets string                      `json:"imagePullSecrets,omitempty"`
	Resources        corev1.ResourceRequirements `json:"resources,omitempty"`
	Env              []corev1.EnvVar             `json:"env"`
	// The bruco stream conf
	Conf BrucoConf `json:"stream"`
}

// BrucoStatus is the status for a Bruco resource
type BrucoStatus struct {
	AvailableReplicas int32 `json:"availableReplicas"`
	CurrentGeneration int64 `json:"currentGeneration"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BrucoList is a list of Bruco resources
type BrucoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Bruco `json:"items"`
}
