// Copyright (c) 2019 Stackinsights to present
// All rights reserved
package v1alpha1

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// BanyanDBSpec defines the desired state of BanyanDB
type BanyanDBSpec struct {
	// Version of BanyanDB.
	// +kubebuilder:validation:Required
	Version string `json:"version"`

	// Number of replicas
	// +kubebuilder:validation:Required
	Counts int `json:"counts"`

	// Pod template of each BanyanDB instance
	// +kubebuilder:validation:Required
	Image string `json:"image"`

	// Pod affinity
	// +kubebuilder:validation:Optional
	Affinity corev1.Affinity `json:"affinity"`

	// BanyanDB startup parameters
	// +kubebuilder:validation:Optional
	Config []string `json:"config"`

	// BanyanDB HTTP Service
	// +kubebuilder:validation:Optional
	HTTPSvc Service `json:"httpService,omitempty"`

	// BanyanDB gRPC Serice
	// +kubebuilder:validation:Optional
	GRPCSvc Service `json:"gRPCService,omitempty"`

	// BanyanDB Storage
	// +kubebuilder:validation:Optional
	Storages []StorageConfig `json:"storages,omitempty"`
}

type StorageConfig struct {
	// the name of storage config
	// +kubebuilder:validation:Optional
	Name string `json:"name,omitempty"`

	// mount path of the volume
	// +kubebuilder:validation:Optional
	Path string `json:"path,omitempty"`

	// the persistent volume spec for the storage
	// +kubebuilder:validation:Optional
	PersistentVolumeClaimSpec corev1.PersistentVolumeClaimSpec `json:"persistentVolumeClaimSpec,omitempty"`
}

// BanyanDBStatus defines the observed state of BanyanDB
type BanyanDBStatus struct {
	AvailableReplicas int32 `json:"available_pods,omitempty"`
	// Represents the latest available observations of the underlying statefulset's current state.
	// +kubebuilder:validation:Optional
	Conditions []appsv1.DeploymentCondition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// BanyanDB is the Schema for the banyandbs API
type BanyanDB struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BanyanDBSpec   `json:"spec,omitempty"`
	Status BanyanDBStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BanyanDBList contains a list of BanyanDB
type BanyanDBList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BanyanDB `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BanyanDB{}, &BanyanDBList{})
}
