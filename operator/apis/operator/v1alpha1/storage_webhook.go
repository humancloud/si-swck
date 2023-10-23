// Copyright (c) 2019 Stackinsights to present
// All rights reserved
package v1alpha1

import (
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var storagelog = logf.Log.WithName("storage-resource")

func (r *Storage) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// nolint: lll
// +kubebuilder:webhook:admissionReviewVersions=v1,sideEffects=None,path=/mutate-operator-stackinsights-apache-org-v1alpha1-storage,mutating=true,failurePolicy=fail,groups=operator.stackinsights.apache.org,resources=storages,verbs=create;update,versions=v1alpha1,name=mstorage.kb.io

var _ webhook.Defaulter = &Storage{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Storage) Default() {
	storagelog.Info("default", "name", r.Name)
	if r.Spec.ConnectType == "internal" {
		if r.Spec.Image == "" {
			r.Spec.Image = "docker.elastic.co/elasticsearch/elasticsearch:7.5.1"
		}
		if r.Spec.Instances == 0 {
			r.Spec.Instances = 3
		}
	}
}

// nolint: lll
// +kubebuilder:webhook:admissionReviewVersions=v1,sideEffects=None,verbs=create;update,path=/validate-operator-stackinsights-apache-org-v1alpha1-storage,mutating=false,failurePolicy=fail,groups=operator.stackinsights.apache.org,resources=storages,versions=v1alpha1,name=vstorage.kb.io

var _ webhook.Validator = &Storage{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Storage) ValidateCreate() error {
	storagelog.Info("validate create", "name", r.Name)
	return r.valid()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Storage) ValidateUpdate(_ runtime.Object) error {
	storagelog.Info("validate update", "name", r.Name)
	return r.valid()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Storage) ValidateDelete() error {
	storagelog.Info("validate delete", "name", r.Name)
	return nil
}

func (r *Storage) valid() error {
	var allErrs field.ErrorList
	if r.Spec.Type != "elasticsearch" {
		storagelog.Info("Invalid Storage Type")
		err := field.Invalid(field.NewPath("spec").Child("type"),
			r.Spec.Type,
			"d. must be elasticsearch")
		allErrs = append(allErrs, err)
	}
	if r.Spec.ConnectType != "internal" && r.Spec.ConnectType != "external" {
		storagelog.Info("Invalid Storage ConnectType")
		err := field.Invalid(field.NewPath("spec").Child("connecttype"),
			r.Spec.ConnectType,
			"d. must be internal or external ")
		allErrs = append(allErrs, err)
	}
	if len(allErrs) != 0 {
		return apierrors.NewInvalid(
			schema.GroupKind{Group: r.GroupVersionKind().Group, Kind: r.GroupVersionKind().Kind},
			r.Name,
			allErrs)
	}
	return nil
}
