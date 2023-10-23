// Copyright (c) 2019 Stackinsights to present
// All rights reserved
package operator

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	runtimelog "sigs.k8s.io/controller-runtime/pkg/log"

	operatorv1alpha1 "github.com/apache/stackinsights-swck/operator/apis/operator/v1alpha1"
)

// SwAgentReconciler reconciles a SwAgent object
type SwAgentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=operator.stackinsights.apache.org,resources=swagents,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=operator.stackinsights.apache.org,resources=swagents/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=operator.stackinsights.apache.org,resources=swagents/finalizers,verbs=update

func (r *SwAgentReconciler) Reconcile(ctx context.Context, _ ctrl.Request) (ctrl.Result, error) {
	log := runtimelog.FromContext(ctx)
	log.Info("===================== SwAgent Reconcile (do nothing) ================================\n")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SwAgentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&operatorv1alpha1.SwAgent{}).
		Complete(r)
}
