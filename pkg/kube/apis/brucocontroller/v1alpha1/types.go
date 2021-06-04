package v1alpha1

import (
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

// BrucoSpec is the spec for a Foo resource
type BrucoSpec struct {
	Replicas       *int32 `json:"replicas"`
	ContainerImage string `json:"containerImage,omitempty"`
	FunctionURL    string `json:"functionURL"`
}

// BrucoStatus is the status for a Foo resource
type BrucoStatus struct {
	AvailableReplicas  int32  `json:"availableReplicas"`
	CurrentFunctionURL string `json:"currentFunctionURL"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BrucoList is a list of Foo resources
type BrucoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Bruco `json:"items"`
}
