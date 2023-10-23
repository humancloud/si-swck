// Copyright (c) 2019 Stackinsights to present
// All rights reserved
package v1alpha1

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var fetcherlog = logf.Log.WithName("fetcher-resource")

func (r *Fetcher) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// nolint: lll
// +kubebuilder:webhook:admissionReviewVersions=v1,sideEffects=None,path=/mutate-operator-stackinsights-apache-org-v1alpha1-fetcher,mutating=true,failurePolicy=fail,groups=operator.stackinsights.apache.org,resources=fetchers,verbs=create;update,versions=v1alpha1,name=mfetcher.kb.io

var _ webhook.Defaulter = &Fetcher{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Fetcher) Default() {
	fetcherlog.Info("default", "name", r.Name)
	if r.Spec.ClusterName == "" {
		r.Spec.ClusterName = r.Name
	}
}

// nolint: lll
// +kubebuilder:webhook:admissionReviewVersions=v1,sideEffects=None,verbs=create;update,path=/validate-operator-stackinsights-apache-org-v1alpha1-fetcher,mutating=false,failurePolicy=fail,groups=operator.stackinsights.apache.org,resources=fetchers,versions=v1alpha1,name=vfetcher.kb.io

var _ webhook.Validator = &Fetcher{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Fetcher) ValidateCreate() error {
	fetcherlog.Info("validate create", "name", r.Name)
	return r.validate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Fetcher) ValidateUpdate(_ runtime.Object) error {
	fetcherlog.Info("validate update", "name", r.Name)
	return r.validate()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Fetcher) ValidateDelete() error {
	fetcherlog.Info("validate delete", "name", r.Name)
	return nil
}

func (r *Fetcher) validate() error {
	if r.Spec.ClusterName == "" {
		return fmt.Errorf("cluster name is absent")
	}
	return nil
}
