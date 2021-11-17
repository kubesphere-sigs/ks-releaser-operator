package controllers

import (
	"fmt"
	devopskubesphereiov1alpha1 "github.com/kubesphere-sigs/ks-releaser-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

func newDeployment(releaser *devopskubesphereiov1alpha1.ReleaserController) *appsv1.Deployment {
	var env []corev1.EnvVar
	if !releaser.Spec.Webhook {
		env = []corev1.EnvVar{{
			Name: "WEBHOOK",
			Value: "false",
		}}
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      releaser.Name,
			Namespace: releaser.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name": "ks-releaser",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app.kubernetes.io/name": "ks-releaser",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app.kubernetes.io/name": "ks-releaser",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "server",
						Image: CombineImageTag(releaser.Spec.Image, releaser.Spec.Version),
						Env: env,
					}},
				},
			},
		},
	}
}

// CombineImageTag will return the combined image and tag in the proper format for tags and digests.
func CombineImageTag(img string, tag string) string {
	if strings.Contains(tag, ":") {
		return fmt.Sprintf("%s@%s", img, tag) // Digest
	} else if len(tag) > 0 {
		return fmt.Sprintf("%s:%s", img, tag) // Tag
	}
	return img // No tag, use default
}
