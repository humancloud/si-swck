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

var banyandbLog = logf.Log.WithName("banyandb-resource")

func (r *BanyanDB) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// nolint: lll
//+kubebuilder:webhook:path=/mutate-operator-stackinsights-apache-org-v1alpha1-banyandb,mutating=true,failurePolicy=fail,sideEffects=None,groups=operator.stackinsights.apache.org,resources=banyandbs,verbs=create;update,versions=v1alpha1,name=mbanyandb.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &BanyanDB{}

func (r *BanyanDB) Default() {
	banyandbLog.Info("default", "name", r.Name)

	if r.Spec.Version == "" {
		// use the latest version by default
		r.Spec.Version = "latest"
	}

	if r.Spec.Image == "" {
		r.Spec.Image = fmt.Sprintf("apache/stackinsights-banyandb:%s", r.Spec.Version)
	}

	if r.Spec.Counts == 0 {
		// currently only support one data copy
		r.Spec.Counts = 1
	}
}

// nolint: lll
// +kubebuilder:webhook:admissionReviewVersions=v1,sideEffects=None,verbs=create;update,path=/validate-operator-stackinsights-apache-org-v1alpha1-banyandb,mutating=false,failurePolicy=fail,groups=operator.stackinsights.apache.org,resources=banyandbs,versions=v1alpha1,name=vbanyandb.kb.io

var _ webhook.Validator = &BanyanDB{}

func (r *BanyanDB) ValidateCreate() error {
	banyandbLog.Info("validate create", "name", r.Name)
	return r.validate()
}

func (r *BanyanDB) ValidateUpdate(_ runtime.Object) error {
	banyandbLog.Info("validate update", "name", r.Name)
	return r.validate()
}

func (r *BanyanDB) ValidateDelete() error {
	banyandbLog.Info("validate delete", "name", r.Name)
	return r.validate()
}

func (r *BanyanDB) validate() error {
	if r.Spec.Image == "" {
		return fmt.Errorf("image is absent")
	}

	if r.Spec.Counts != 1 {
		return fmt.Errorf("banyandb only support 1 copy for now")
	}

	return nil
}
