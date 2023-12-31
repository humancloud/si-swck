// Copyright (c) 2019 Stackinsights to present
// All rights reserved
package operator

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	apiequal "k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/retry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	runtimelog "sigs.k8s.io/controller-runtime/pkg/log"

	operatorv1alpha1 "github.com/humancloud/si-swck/operator/apis/operator/v1alpha1"
	"github.com/humancloud/si-swck/operator/pkg/kubernetes"
)

// SatelliteReconciler reconciles a Satellite object
type SatelliteReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	FileRepo kubernetes.Repo
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=operator.stackinsights.apache.org,resources=satellites,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=operator.stackinsights.apache.org,resources=satellites/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=operator.stackinsights.apache.org,resources=satellites/finalizers,verbs=update

func (r *SatelliteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := runtimelog.FromContext(ctx)
	log.Info("=====================satellite reconcile started================================")

	satellite := operatorv1alpha1.Satellite{}
	if err := r.Client.Get(ctx, req.NamespacedName, &satellite); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	ff, err := r.FileRepo.GetFilesRecursive("templates")
	if err != nil {
		log.Error(err, "failed to load resource templates")
		return ctrl.Result{}, err
	}
	app := kubernetes.Application{
		Client:   r.Client,
		FileRepo: r.FileRepo,
		CR:       &satellite,
		GVK:      operatorv1alpha1.GroupVersion.WithKind("Satellite"),
		Recorder: r.Recorder,
	}

	if err := app.ApplyAll(ctx, ff, log); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.checkState(ctx, log, &satellite); err != nil {
		log.Error(err, "failed to check sub resources state")
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: schedDuration}, nil
}

func (r *SatelliteReconciler) checkState(ctx context.Context, log logr.Logger, satellite *operatorv1alpha1.Satellite) error {
	overlay := operatorv1alpha1.SatelliteStatus{}
	deployment := apps.Deployment{}
	errCol := new(kubernetes.ErrorCollector)
	if err := r.Client.Get(ctx, client.ObjectKey{Namespace: satellite.Namespace, Name: satellite.Name + "-satellite"}, &deployment); err != nil {
		errCol.Collect(fmt.Errorf("failed to get deployment :%v", err))
	} else {
		overlay.Conditions = deployment.Status.Conditions
		overlay.AvailableReplicas = deployment.Status.AvailableReplicas
	}
	service := core.Service{}
	if err := r.Client.Get(ctx, client.ObjectKey{Namespace: satellite.Namespace, Name: satellite.Name + "-satellite"}, &service); err != nil {
		errCol.Collect(fmt.Errorf("failed to get service :%v", err))
	} else {
		overlay.Address = fmt.Sprintf("%s.%s", service.Name, service.Namespace)
	}
	if apiequal.Semantic.DeepDerivative(overlay, satellite.Status) {
		log.Info("Status keeps the same as before")
		return errCol.Error()
	}
	satellite.Status = overlay
	satellite.Kind = "Satellite"
	if err := kubernetes.ApplyOverlay(satellite, &operatorv1alpha1.Satellite{Status: overlay}); err != nil {
		errCol.Collect(fmt.Errorf("failed to apply overlay: %w", err))
		return errCol.Error()
	}
	if err := r.updateStatus(ctx, satellite, overlay, errCol); err != nil {
		errCol.Collect(fmt.Errorf("failed to update status of satellite: %w", err))
	}
	log.Info("updated Status sub resource")

	return errCol.Error()
}

func (r *SatelliteReconciler) updateStatus(ctx context.Context, satellite *operatorv1alpha1.Satellite,
	overlay operatorv1alpha1.SatelliteStatus, errCol *kubernetes.ErrorCollector) error {
	// avoid resource conflict
	return retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		if err := r.Client.Get(ctx, client.ObjectKey{Name: satellite.Name, Namespace: satellite.Namespace}, satellite); err != nil {
			errCol.Collect(fmt.Errorf("failed to get satellite: %w", err))
		}
		satellite.Status = overlay
		satellite.Kind = "Satellite"
		if err := kubernetes.ApplyOverlay(satellite, &operatorv1alpha1.Satellite{Status: overlay}); err != nil {
			errCol.Collect(fmt.Errorf("failed to apply overlay: %w", err))
		}
		if err := r.Status().Update(ctx, satellite); err != nil {
			errCol.Collect(fmt.Errorf("failed to update status of satellite: %w", err))
		}
		return errCol.Error()
	})
}

// SetupWithManager sets up the controller with the Manager.
func (r *SatelliteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&operatorv1alpha1.Satellite{}).
		Complete(r)
}
