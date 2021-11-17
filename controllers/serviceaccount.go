package controllers

import (
	"context"
	"fmt"
	devopskubesphereiov1alpha1 "github.com/kubesphere-sigs/ks-releaser-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReleaserControllerReconciler) reconcileServiceAccount(name string, cr *devopskubesphereiov1alpha1.ReleaserController) (
	sa *corev1.ServiceAccount, err error) {
	sa = newServiceAccountWithName(name, cr)
	if err = r.Client.Get(context.TODO(), types.NamespacedName{
		Namespace: sa.Namespace,
		Name:      sa.Name,
	}, sa); err != nil {
		if !errors.IsNotFound(err) {
			return
		}
		err = r.Client.Create(context.TODO(), sa)
	}

	if err == nil {
		err = controllerutil.SetControllerReference(cr, sa, r.Scheme)
	}
	return
}

func newServiceAccountWithName(name string, cr *devopskubesphereiov1alpha1.ReleaserController) *corev1.ServiceAccount {
	sa := newServiceAccount(cr)
	sa.ObjectMeta.Name = fmt.Sprintf("%s-%s", cr.Name, name)
	return sa
}

// newServiceAccount returns a new ServiceAccount instance.
func newServiceAccount(cr *devopskubesphereiov1alpha1.ReleaserController) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
		},
	}
}
