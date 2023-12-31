// Copyright (c) 2019 Stackinsights to present
// All rights reserved
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// JavaAgentSpec defines the desired state of JavaAgent
type JavaAgentSpec struct {
	// PodSelector is the selector label of injected Pod
	// +kubebuilder:validation:Optional
	PodSelector string `json:"podSelector,omitempty"`
	// ServiceName is the name of service in the injected agent, which need to be printed
	// +kubebuilder:validation:Required
	ServiceName string `json:"serviceName,omitempty"`
	// BackendService is the backend service in the injected agent, which need to be printed
	// +kubebuilder:validation:Required
	BackendService string `json:"backendService,omitempty"`
	// AgentConfiguration is the injected agent's final configuration
	// +kubebuilder:validation:Optional
	AgentConfiguration map[string]string `json:"agentConfiguration,omitempty"`
}

// JavaAgentStatus defines the observed state of JavaAgent
type JavaAgentStatus struct {
	// The number of pods that need to be injected
	ExpectedInjectedNum int `json:"expectedInjectiedNum,omitempty"`
	// The number of pods that injected successfully
	RealInjectedNum int `json:"realInjectedNum,omitempty"`
	// The time the JavaAgent was created.
	CreationTime metav1.Time `json:"creationTime,omitempty"`
	// The last time this condition was updated.
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="PodSelector",type="string",JSONPath=".spec.podSelector",description="The selector label of injected Pod"
// +kubebuilder:printcolumn:name="ServiceName",type="string",JSONPath=".spec.serviceName",description="The name of service in the injected agent"
// +kubebuilder:printcolumn:name="BackendService",type="string",JSONPath=".spec.backendService",description="The backend service in the injected agent"

// JavaAgent is the Schema for the javaagents API
type JavaAgent struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JavaAgentSpec   `json:"spec,omitempty"`
	Status JavaAgentStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// JavaAgentList contains a list of JavaAgent
type JavaAgentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []JavaAgent `json:"items"`
}

func init() {
	SchemeBuilder.Register(&JavaAgent{}, &JavaAgentList{})
}
