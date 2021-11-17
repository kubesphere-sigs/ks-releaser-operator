package controllers

import (
	"context"
	"fmt"
	devopskubesphereiov1alpha1 "github.com/kubesphere-sigs/ks-releaser-operator/api/v1alpha1"
	v1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func newRole(name string, rules []v1.PolicyRule, cr *devopskubesphereiov1alpha1.ReleaserController) *v1.Role {
	return &v1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      generateResourceName(name, cr),
			Namespace: cr.Namespace,
		},
		Rules: rules,
	}
}

func newClusterRole(name string, rules []v1.PolicyRule, cr *devopskubesphereiov1alpha1.ReleaserController) *v1.ClusterRole {
	return &v1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:      generateResourceName(name, cr),
			Namespace: cr.Namespace,
		},
		Rules: rules,
	}
}

func generateResourceName(componentName string, cr *devopskubesphereiov1alpha1.ReleaserController) string {
	return cr.Name + "-" + componentName
}

func (r *ReleaserControllerReconciler) reconcileRole(name string, rules []v1.PolicyRule, cr *devopskubesphereiov1alpha1.ReleaserController) (
	role *v1.Role, err error) {
	role = newRole(name, rules, cr)
	existingRole := &v1.Role{}
	if err = r.Client.Get(context.TODO(), types.NamespacedName{Namespace: role.Namespace, Name: role.Name}, existingRole); err != nil {
		if !errors.IsNotFound(err) {
			return nil, fmt.Errorf("failed to reconcile the role: %v", err)
		}

		err = r.Client.Create(context.TODO(), role)
		return
	}

	existingRole.Rules = role.Rules
	err = r.Client.Update(context.TODO(), existingRole)
	return
}

func (r *ReleaserControllerReconciler) reconcileClusterRole(name string, rules []v1.PolicyRule,
	cr *devopskubesphereiov1alpha1.ReleaserController) (role *v1.ClusterRole, err error) {
	role = newClusterRole(name, rules, cr)
	existingRole := &v1.ClusterRole{}
	if err = r.Client.Get(context.TODO(), types.NamespacedName{Namespace: role.Namespace, Name: role.Name}, existingRole); err != nil {
		if !errors.IsNotFound(err) {
			return nil, fmt.Errorf("failed to reconcile the role: %v", err)
		}

		err = r.Client.Create(context.TODO(), role)
		return
	}

	existingRole.Rules = role.Rules
	err = r.Client.Update(context.TODO(), existingRole)
	return
}
