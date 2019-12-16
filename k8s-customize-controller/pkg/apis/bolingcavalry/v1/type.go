package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Programmer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ProgrammerSpec `json:"spec"`
	Status ProgrammerStatus `json:"status"`
}

type ProgrammerSpec struct {
	name   string `json:"name"`
	company string `json:"company"`
	DeploymentName string `json:"deploymentName"`
	Replicas *int32 `json:"replicas"`
}

type ProgrammerStatus struct {
	AvailableReplicas int32 `json:"availableReplicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ProgrammerList is a list of Programmer resources
type ProgrammerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Programmer `json:"items"`
}
