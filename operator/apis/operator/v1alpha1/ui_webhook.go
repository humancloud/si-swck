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
var uilog = logf.Log.WithName("ui-resource")

func (r *UI) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// nolint: lll
// +kubebuilder:webhook:admissionReviewVersions=v1,sideEffects=None,path=/mutate-operator-stackinsights-apache-org-v1alpha1-ui,mutating=true,failurePolicy=fail,groups=operator.stackinsights.apache.org,resources=uis,verbs=create;update,versions=v1alpha1,name=mui.kb.io

var _ webhook.Defaulter = &UI{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *UI) Default() {
	uilog.Info("default", "name", r.Name)

	if r.Spec.Image == "" {
		r.Spec.Image = fmt.Sprintf("apache/stackinsights-ui:%s", r.Spec.Version)
	}

	r.Spec.Service.Template.Default()
	if r.Spec.OAPServerAddress == "" {
		r.Spec.OAPServerAddress = fmt.Sprintf("http://%s-oap.%s:12800", r.Name, r.Namespace)
	}
}

// nolint: lll
// +kubebuilder:webhook:admissionReviewVersions=v1,sideEffects=None,verbs=create;update,path=/validate-operator-stackinsights-apache-org-v1alpha1-ui,mutating=false,failurePolicy=fail,groups=operator.stackinsights.apache.org,resources=uis,versions=v1alpha1,name=vui.kb.io

var _ webhook.Validator = &UI{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *UI) ValidateCreate() error {
	uilog.Info("validate create", "name", r.Name)
	return r.validate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *UI) ValidateUpdate(_ runtime.Object) error {
	uilog.Info("validate update", "name", r.Name)
	return r.validate()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *UI) ValidateDelete() error {
	uilog.Info("validate delete", "name", r.Name)
	return nil
}

func (r *UI) validate() error {
	if r.Spec.Image == "" {
		return fmt.Errorf("image is absent")
	}
	if err := r.Spec.Service.Template.Validate(); err != nil {
		return fmt.Errorf("service template is invalid: %w", err)
	}
	if r.Spec.OAPServerAddress == "" {
		return fmt.Errorf("oap server address is absent")
	}
	return nil
}
