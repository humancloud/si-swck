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
var satellitelog = logf.Log.WithName("satellite-resource")

func (r *Satellite) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//nolint: lll
//+kubebuilder:webhook:path=/mutate-operator-stackinsights-apache-org-v1alpha1-satellite,mutating=true,failurePolicy=fail,sideEffects=None,groups=operator.stackinsights.apache.org,resources=satellites,verbs=create;update,versions=v1alpha1,name=msatellite.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Satellite{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Satellite) Default() {
	satellitelog.Info("default", "name", r.Name)

	image := r.Spec.Image
	if image == "" {
		r.Spec.Image = fmt.Sprintf("apache/stackinsights-satellite:%s", r.Spec.Version)
	}

	oapServerName := r.Spec.OAPServerName
	if oapServerName == "" {
		r.Spec.OAPServerName = r.Name
	}

	if r.ObjectMeta.Annotations[annotationKeyIstioSetup] == "" {
		r.Annotations[annotationKeyIstioSetup] = fmt.Sprintf("istioctl install --set profile=demo "+
			"--set meshConfig.defaultConfig.envoyAccessLogService.address=%s.%s:11800 "+
			"--set meshConfig.enableEnvoyAccessLogService=true", r.Name, r.Namespace)
	}
}

//nolint: lll
//+kubebuilder:webhook:path=/validate-operator-stackinsights-apache-org-v1alpha1-satellite,mutating=false,failurePolicy=fail,sideEffects=None,groups=operator.stackinsights.apache.org,resources=satellites,verbs=create;update,versions=v1alpha1,name=vsatellite.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Satellite{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Satellite) ValidateCreate() error {
	satellitelog.Info("validate create", "name", r.Name)
	return r.validate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Satellite) ValidateUpdate(_ runtime.Object) error {
	satellitelog.Info("validate update", "name", r.Name)
	return r.validate()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Satellite) ValidateDelete() error {
	satellitelog.Info("validate delete", "name", r.Name)
	return r.validate()
}

func (r *Satellite) validate() error {
	if r.Spec.Image == "" {
		return fmt.Errorf("image is absent")
	}
	return nil
}
