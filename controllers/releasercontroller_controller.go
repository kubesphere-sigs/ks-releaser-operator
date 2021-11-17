/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"github.com/kubesphere-sigs/ks-releaser-operator/common"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	devopskubesphereiov1alpha1 "github.com/kubesphere-sigs/ks-releaser-operator/api/v1alpha1"
)

// ReleaserControllerReconciler reconciles a ReleaserController object
type ReleaserControllerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=devops.kubesphere.io,resources=releasercontrollers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=devops.kubesphere.io,resources=releasercontrollers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=devops.kubesphere.io,resources=releasercontrollers/finalizers,verbs=update
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterroles;clusterrolebindings;roles;rolebindings,verbs=*
//+kubebuilder:rbac:groups="",resources=configmaps;endpoints;events;pods;namespaces;secrets;serviceaccounts;services;services/finalizers,verbs=*
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ReleaserController object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *ReleaserControllerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	releaserCtl := &devopskubesphereiov1alpha1.ReleaserController{}
	err := r.Client.Get(ctx, req.NamespacedName, releaserCtl)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	if releaserCtl.GetDeletionTimestamp() != nil {
		if releaserCtl.IsDeletionFinalizerPresent() {
			if err = r.deleteClusterResources(releaserCtl); err != nil {
				return reconcile.Result{}, err
			}

			if err = r.removeDeletionFinalizer(releaserCtl); err != nil {
				return reconcile.Result{}, err
			}
		}
		return reconcile.Result{}, nil
	}

	if !releaserCtl.IsDeletionFinalizerPresent() {
		if err = r.addDeletionFinalizer(releaserCtl); err != nil {
			return reconcile.Result{}, err
		}
	}

	if err = r.Client.Get(ctx, req.NamespacedName, releaserCtl); err != nil {
		return reconcile.Result{}, err
	}

	if err := r.reconcileResources(releaserCtl); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *ReleaserControllerReconciler) deleteClusterResources(cr *devopskubesphereiov1alpha1.ReleaserController) (err error) {
	existingDeploy := &appsv1.Deployment{}

	if err = r.Client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, existingDeploy); err != nil {
		if !errors.IsNotFound(err) {
			return
		}
	} else {
		if err = r.Client.Delete(context.TODO(), existingDeploy); err != nil {
			return
		}
	}

	existingServiceAccount := &corev1.ServiceAccount{}
	if err = r.Client.Get(context.TODO(), types.NamespacedName{
		Namespace: cr.Namespace,
		Name:      fmt.Sprintf("%s-%s", cr.Name, component),
	}, existingServiceAccount); err != nil {
		if !errors.IsNotFound(err) {
			return
		}
	} else {
		if err = r.Client.Delete(context.TODO(), existingServiceAccount); err != nil {
			return
		}
	}

	existingRole := &v1.Role{}
	if err = r.Client.Get(context.TODO(), types.NamespacedName{
		Namespace: cr.Namespace,
		Name:      generateResourceName(component, cr),
	}, existingRole); err != nil {
		if !errors.IsNotFound(err) {
			return
		}
	} else {
		if err = r.Client.Delete(context.TODO(), existingRole); err != nil {
			return
		}
	}

	existingRoleBinding := &v1.RoleBinding{}
	if err = r.Client.Get(context.TODO(), types.NamespacedName{
		Namespace: cr.Namespace,
		Name:      generateResourceName(component, cr),
	}, existingRoleBinding); err != nil {
		if !errors.IsNotFound(err) {
			return
		}
	} else {
		if err = r.Client.Delete(context.TODO(), existingRoleBinding); err != nil {
			return
		}
	}

	existingClusterRole := &v1.ClusterRole{}
	if err = r.Client.Get(context.TODO(), types.NamespacedName{
		Namespace: cr.Namespace,
		Name:      generateResourceName(component, cr),
	}, existingClusterRole); err != nil {
		if !errors.IsNotFound(err) {
			return
		}
	} else {
		if err = r.Client.Delete(context.TODO(), existingClusterRole); err != nil {
			return
		}
	}

	existingClusterRoleBinding := &v1.ClusterRoleBinding{}
	if err = r.Client.Get(context.TODO(), types.NamespacedName{
		Namespace: cr.Namespace,
		Name:      generateResourceName(component, cr),
	}, existingClusterRoleBinding); err != nil {
		if !errors.IsNotFound(err) {
			return
		}
	} else {
		if err = r.Client.Delete(context.TODO(), existingClusterRoleBinding); err != nil {
			return
		}
	}
	return
}

func instanceSelector(name string) (labels.Selector, error) {
	selector := labels.NewSelector()
	requirement, err := labels.NewRequirement(common.ManagedByLabel, selection.Equals, []string{name})
	if err != nil {
		return nil, fmt.Errorf("failed to create a requirement for %w", err)
	}
	return selector.Add(*requirement), nil
}

func (r *ReleaserControllerReconciler) removeDeletionFinalizer(cr *devopskubesphereiov1alpha1.ReleaserController) error {
	cr.Finalizers = removeString(cr.GetFinalizers(), common.DeletionFinalizer)
	if err := r.Client.Update(context.TODO(), cr); err != nil {
		return fmt.Errorf("failed to remove deletion finalizer from %s: %w", cr.Name, err)
	}
	return nil
}

func removeString(slice []string, s string) []string {
	var result []string
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return result
}

func (r *ReleaserControllerReconciler) addDeletionFinalizer(cr *devopskubesphereiov1alpha1.ReleaserController) error {
	cr.Finalizers = append(cr.Finalizers, common.DeletionFinalizer)
	if err := r.Client.Update(context.TODO(), cr); err != nil {
		return fmt.Errorf("failed to add deletion finalizer for %s: %w", cr.Name, err)
	}
	return nil
}

func (r *ReleaserControllerReconciler) reconcileResources(cr *devopskubesphereiov1alpha1.ReleaserController) (err error) {
	if err = r.reconcileRoles(cr); err != nil {
		return
	}

	if _, err = r.reconcileServiceAccounts(cr); err != nil {
		return
	}

	if err = r.reconcileDeployment(cr); err != nil {
		return
	}
	return
}

func (r *ReleaserControllerReconciler) reconcileDeployment(cr *devopskubesphereiov1alpha1.ReleaserController) (err error) {
	deploy := newDeployment(cr)
	deploy.Spec.Template.Spec.ServiceAccountName = fmt.Sprintf("%s-%s", cr.Name, component)

	existingDeploy := &appsv1.Deployment{}
	if err = r.Client.Get(context.TODO(), types.NamespacedName{
		Namespace: deploy.Namespace,
		Name:      deploy.Name,
	}, existingDeploy); err != nil {
		if !errors.IsNotFound(err) {
			return
		}

		if err = r.Client.Create(context.TODO(), deploy); err != nil {
			return
		}
	} else {
		existingDeploy.Spec = deploy.Spec
		err = r.Client.Update(context.TODO(), existingDeploy)
	}
	return
}

const component = "ks-releaser-controller"

func (r *ReleaserControllerReconciler) reconcileRoles(cr *devopskubesphereiov1alpha1.ReleaserController) (err error) {
	if _, err = r.reconcileRole(component, policyRuleForLeaderElection(), cr); err != nil {
		return
	}

	if err = r.reconcileRoleBindings(component, policyRuleForLeaderElection(), cr); err != nil {
		return
	}

	if _, err = r.reconcileClusterRole(component, policyRuleForManager(), cr); err != nil {
		return
	}

	if err = r.reconcileClusterRoleBindings(component, policyRuleForManager(), cr); err != nil {
		return
	}
	return
}

func (r *ReleaserControllerReconciler) reconcileRoleBindings(name string, rules []v1.PolicyRule, cr *devopskubesphereiov1alpha1.ReleaserController) (err error) {
	var sa *corev1.ServiceAccount
	var role *v1.Role

	if sa, err = r.reconcileServiceAccount(name, cr); err != nil {
		return
	}

	if role, err = r.reconcileRole(name, rules, cr); err != nil {
		return
	}

	roleBinding := newRoleBindingWithName(name, cr)
	roleBinding.Namespace = role.Namespace

	existingRoleBinding := &v1.RoleBinding{}
	if err = r.Client.Get(context.TODO(), types.NamespacedName{Name: roleBinding.Name, Namespace: roleBinding.Namespace}, existingRoleBinding); err != nil {
		if !errors.IsNotFound(err) {
			return fmt.Errorf("failed to get rolebinding, error: %v", err)
		}

		roleBinding.Subjects = []v1.Subject{{
			Kind:      v1.ServiceAccountKind,
			Name:      sa.Name,
			Namespace: sa.Namespace,
		}}
		roleBinding.RoleRef = v1.RoleRef{
			APIGroup: v1.GroupName,
			Kind:     "Role",
			Name:     role.Name,
		}
		err = r.Client.Create(context.TODO(), roleBinding)
	}
	return
}

func (r *ReleaserControllerReconciler) reconcileClusterRoleBindings(name string, rules []v1.PolicyRule, cr *devopskubesphereiov1alpha1.ReleaserController) (err error) {
	var sa *corev1.ServiceAccount
	var role *v1.ClusterRole

	if sa, err = r.reconcileServiceAccount(name, cr); err != nil {
		return
	}

	if role, err = r.reconcileClusterRole(name, rules, cr); err != nil {
		return
	}

	roleBinding := newClusterRoleBindingWithName(name, cr)

	existingRoleBinding := &v1.ClusterRoleBinding{}
	if err = r.Client.Get(context.TODO(), types.NamespacedName{Name: roleBinding.Name}, existingRoleBinding); err != nil {
		if !errors.IsNotFound(err) {
			return fmt.Errorf("failed to get clusterrolebinding, error: %v", err)
		}

		roleBinding.Subjects = []v1.Subject{{
			Kind:      v1.ServiceAccountKind,
			Name:      sa.Name,
			Namespace: sa.Namespace,
		}}
		roleBinding.RoleRef = v1.RoleRef{
			APIGroup: v1.GroupName,
			Kind:     "ClusterRole",
			Name:     role.Name,
		}
		err = r.Client.Create(context.TODO(), roleBinding)
	}
	return
}

func (r *ReleaserControllerReconciler) reconcileServiceAccounts(cr *devopskubesphereiov1alpha1.ReleaserController) (sa *corev1.ServiceAccount, err error) {
	if sa, err = r.reconcileServiceAccount(component, cr); err != nil {
		return
	}
	return
}

func (r *ReleaserControllerReconciler) reconcileServices(cr *devopskubesphereiov1alpha1.ReleaserController) (err error) {
	return
}

func (r *ReleaserControllerReconciler) reconcileConfigMaps(cr *devopskubesphereiov1alpha1.ReleaserController) (err error) {
	return
}

// SetupWithManager sets up the controller with the Manager.
func (r *ReleaserControllerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&devopskubesphereiov1alpha1.ReleaserController{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
