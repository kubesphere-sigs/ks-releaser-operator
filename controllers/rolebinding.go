package controllers

import (
	"fmt"
	devopskubesphereiov1alpha1 "github.com/kubesphere-sigs/ks-releaser-operator/api/v1alpha1"
	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newRoleBindingWithName(name string, cr *devopskubesphereiov1alpha1.ReleaserController) *v1.RoleBinding {
	roleBinding := newRoleBinding(cr)
	roleBinding.ObjectMeta.Name = fmt.Sprintf("%s-%s", cr.Name, name)
	return roleBinding
}

func newRoleBinding(cr *devopskubesphereiov1alpha1.ReleaserController) *v1.RoleBinding {
	return &v1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
		},
	}
}

func newClusterRoleBindingWithName(name string, cr *devopskubesphereiov1alpha1.ReleaserController) *v1.ClusterRoleBinding {
	roleBinding := newClusterRoleBinding(cr)
	roleBinding.ObjectMeta.Name = fmt.Sprintf("%s-%s", cr.Name, name)
	return roleBinding
}

func newClusterRoleBinding(cr *devopskubesphereiov1alpha1.ReleaserController) *v1.ClusterRoleBinding {
	return &v1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: cr.Name,
		},
	}
}
