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

const annotationKeyIstioSetup = "istio-setup-command"

// log is for logging in this package.
var oapserverlog = logf.Log.WithName("oapserver-resource")

func (r *OAPServer) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// nolint: lll
// +kubebuilder:webhook:admissionReviewVersions=v1,sideEffects=None,path=/mutate-operator-stackinsights-apache-org-v1alpha1-oapserver,mutating=true,failurePolicy=fail,groups=operator.stackinsights.apache.org,resources=oapservers,verbs=create;update,versions=v1alpha1,name=moapserver.kb.io

var _ webhook.Defaulter = &OAPServer{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *OAPServer) Default() {
	oapserverlog.Info("default", "name", r.Name)

	image := r.Spec.Image
	if image == "" {
		r.Spec.Image = fmt.Sprintf("apache/stackinsights-oap-server:%s", r.Spec.Version)
	}
	for _, envVar := range r.Spec.Config {
		if envVar.Name == "SW_ENVOY_METRIC_ALS_HTTP_ANALYSIS" &&
			r.ObjectMeta.Annotations[annotationKeyIstioSetup] == "" {
			r.Annotations[annotationKeyIstioSetup] = fmt.Sprintf("istioctl install --set profile=demo "+
				"--set meshConfig.defaultConfig.envoyAccessLogService.address=%s.%s:11800 "+
				"--set meshConfig.enableEnvoyAccessLogService=true", r.Name, r.Namespace)
		}
	}
}

// nolint: lll
// +kubebuilder:webhook:admissionReviewVersions=v1,sideEffects=None,verbs=create;update,path=/validate-operator-stackinsights-apache-org-v1alpha1-oapserver,mutating=false,failurePolicy=fail,groups=operator.stackinsights.apache.org,resources=oapservers,versions=v1alpha1,name=voapserver.kb.io

var _ webhook.Validator = &OAPServer{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *OAPServer) ValidateCreate() error {
	oapserverlog.Info("validate create", "name", r.Name)
	return r.validate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *OAPServer) ValidateUpdate(_ runtime.Object) error {
	oapserverlog.Info("validate update", "name", r.Name)
	return r.validate()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *OAPServer) ValidateDelete() error {
	oapserverlog.Info("validate delete", "name", r.Name)
	return nil
}

func (r *OAPServer) validate() error {
	if r.Spec.Image == "" {
		return fmt.Errorf("image is absent")
	}
	return nil
}
