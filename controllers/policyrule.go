package controllers

import v1 "k8s.io/api/rbac/v1"

func policyRuleForLeaderElection() []v1.PolicyRule {
	return []v1.PolicyRule{{
		APIGroups: []string{""},
		Resources: []string{"configmaps"},
		Verbs:     []string{"get", "list", "watch", "create", "update", "patch", "delete"},
	}, {
		APIGroups: []string{"coordination.k8s.io"},
		Resources: []string{"leases"},
		Verbs:     []string{"get", "list", "watch", "create", "update", "patch", "delete"},
	}, {
		APIGroups: []string{""},
		Resources: []string{"events"},
		Verbs:     []string{"create", "patch"},
	}}
}

func policyRuleForManager() []v1.PolicyRule {
	return []v1.PolicyRule{{
		APIGroups: []string{""},
		Resources: []string{"secrets"},
		Verbs:     []string{"get", "list", "watch"},
	}, {
		APIGroups: []string{"devops.kubesphere.io"},
		Resources: []string{"releasers"},
		Verbs:     []string{"create", "delete", "get", "list", "patch", "update", "watch"},
	}, {
		APIGroups: []string{"devops.kubesphere.io"},
		Resources: []string{"releasers/finalizers"},
		Verbs:     []string{"update"},
	}, {
		APIGroups: []string{"devops.kubesphere.io"},
		Resources: []string{"releasers/status"},
		Verbs:     []string{"get", "patch", "update"},
	}}
}
