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
var oapserverdynamicconfiglog = logf.Log.WithName("oapserverdynamicconfig-resource")

func (r *OAPServerDynamicConfig) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// nolint: lll
//+kubebuilder:webhook:path=/mutate-operator-stackinsights-apache-org-v1alpha1-oapserverdynamicconfig,mutating=true,failurePolicy=fail,sideEffects=None,groups=operator.stackinsights.apache.org,resources=oapserverdynamicconfigs,verbs=create;update,versions=v1alpha1,name=moapserverdynamicconfig.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &OAPServerDynamicConfig{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *OAPServerDynamicConfig) Default() {
	oapserverdynamicconfiglog.Info("default", "name", r.Name)

	// Default version is "9.5.0"
	if r.Spec.Version == "" {
		r.Spec.Version = "9.5.0"
	}

	// Default labelselector is "app=collector,release=stackinsights"
	if r.Spec.LabelSelector == "" {
		r.Spec.LabelSelector = "app=collector,release=stackinsights"
	}
}

// nolint: lll
//+kubebuilder:webhook:path=/validate-operator-stackinsights-apache-org-v1alpha1-oapserverdynamicconfig,mutating=false,failurePolicy=fail,sideEffects=None,groups=operator.stackinsights.apache.org,resources=oapserverdynamicconfigs,verbs=create;update,versions=v1alpha1,name=voapserverdynamicconfig.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &OAPServerDynamicConfig{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *OAPServerDynamicConfig) ValidateCreate() error {
	oapserverdynamicconfiglog.Info("validate create", "name", r.Name)
	return r.validate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *OAPServerDynamicConfig) ValidateUpdate(_ runtime.Object) error {
	oapserverdynamicconfiglog.Info("validate update", "name", r.Name)
	return r.validate()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *OAPServerDynamicConfig) ValidateDelete() error {
	oapserverdynamicconfiglog.Info("validate delete", "name", r.Name)
	return nil
}

func (r *OAPServerDynamicConfig) validate() error {
	if r.Spec.Version == "" {
		return fmt.Errorf("OAPServerDynamicConfig's version is absent")
	} else if r.Spec.LabelSelector == "" {
		return fmt.Errorf("OAPServerDynamicConfig's labelselector is absent")
	}
	return nil
}
