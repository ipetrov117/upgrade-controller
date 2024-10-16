package upgrade

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
)

func TestKubernetesControlPlanePlan_RKE2(t *testing.T) {
	version := "v1.30.2+rke2r1"
	addLabels := map[string]string{
		"lifecycle.suse.com/x": "z",
	}

	expectedLabels := map[string]string{
		"lifecycle.suse.com/x": "z",
		"k8s-upgrade":          "control-plane",
	}

	upgradePlan := KubernetesControlPlanePlan(planNameSuffix, version, false, addLabels)
	require.NotNil(t, upgradePlan)

	assert.Equal(t, "Plan", upgradePlan.TypeMeta.Kind)
	assert.Equal(t, "upgrade.cattle.io/v1", upgradePlan.TypeMeta.APIVersion)

	assert.Equal(t, "control-plane-v1-30-2-rke2r1-abcdef", upgradePlan.ObjectMeta.Name)
	assert.Equal(t, "cattle-system", upgradePlan.ObjectMeta.Namespace)
	assert.Equal(t, expectedLabels, upgradePlan.ObjectMeta.Labels)

	require.Len(t, upgradePlan.Spec.NodeSelector.MatchLabels, 0)
	require.Len(t, upgradePlan.Spec.NodeSelector.MatchExpressions, 1)

	matchExpression := upgradePlan.Spec.NodeSelector.MatchExpressions[0]
	assert.Equal(t, "node-role.kubernetes.io/control-plane", matchExpression.Key)
	assert.EqualValues(t, "In", matchExpression.Operator)
	assert.Equal(t, []string{"true"}, matchExpression.Values)

	require.Nil(t, upgradePlan.Spec.Prepare)

	require.NotNil(t, upgradePlan.Spec.Upgrade)
	assert.Equal(t, "rancher/rke2-upgrade", upgradePlan.Spec.Upgrade.Image)
	assert.Nil(t, upgradePlan.Spec.Upgrade.Args)

	assert.Equal(t, version, upgradePlan.Spec.Version)
	assert.EqualValues(t, 1, upgradePlan.Spec.Concurrency)
	assert.True(t, upgradePlan.Spec.Cordon)

	assert.Equal(t, "system-upgrade-controller", upgradePlan.Spec.ServiceAccountName)
	assert.Nil(t, upgradePlan.Spec.Drain)

	tolerations := []corev1.Toleration{
		{
			Key:      "CriticalAddonsOnly",
			Operator: "Equal",
			Value:    "true",
			Effect:   "NoExecute",
		},
		{
			Key:      ControlPlaneLabel,
			Operator: "Equal",
			Value:    "",
			Effect:   "NoSchedule",
		},
		{
			Key:      "node-role.kubernetes.io/etcd",
			Operator: "Equal",
			Value:    "",
			Effect:   "NoExecute",
		},
	}
	assert.Equal(t, tolerations, upgradePlan.Spec.Tolerations)
}

func TestKubernetesControlPlanePlan_K3s(t *testing.T) {
	version := "v1.30.2+k3s1"
	addLabels := map[string]string{
		"lifecycle.suse.com/x": "z",
	}

	expectedLabels := map[string]string{
		"lifecycle.suse.com/x": "z",
		"k8s-upgrade":          "control-plane",
	}

	upgradePlan := KubernetesControlPlanePlan(planNameSuffix, version, false, addLabels)
	require.NotNil(t, upgradePlan)

	assert.Equal(t, "Plan", upgradePlan.TypeMeta.Kind)
	assert.Equal(t, "upgrade.cattle.io/v1", upgradePlan.TypeMeta.APIVersion)

	assert.Equal(t, "control-plane-v1-30-2-k3s1-abcdef", upgradePlan.ObjectMeta.Name)
	assert.Equal(t, "cattle-system", upgradePlan.ObjectMeta.Namespace)
	assert.Equal(t, expectedLabels, upgradePlan.ObjectMeta.Labels)

	require.Len(t, upgradePlan.Spec.NodeSelector.MatchLabels, 0)
	require.Len(t, upgradePlan.Spec.NodeSelector.MatchExpressions, 1)

	matchExpression := upgradePlan.Spec.NodeSelector.MatchExpressions[0]
	assert.Equal(t, "node-role.kubernetes.io/control-plane", matchExpression.Key)
	assert.EqualValues(t, "In", matchExpression.Operator)
	assert.Equal(t, []string{"true"}, matchExpression.Values)

	require.Nil(t, upgradePlan.Spec.Prepare)

	require.NotNil(t, upgradePlan.Spec.Upgrade)
	assert.Equal(t, "rancher/k3s-upgrade", upgradePlan.Spec.Upgrade.Image)
	assert.Nil(t, upgradePlan.Spec.Upgrade.Args)

	assert.Equal(t, version, upgradePlan.Spec.Version)
	assert.EqualValues(t, 1, upgradePlan.Spec.Concurrency)
	assert.True(t, upgradePlan.Spec.Cordon)

	assert.Equal(t, "system-upgrade-controller", upgradePlan.Spec.ServiceAccountName)
	assert.Nil(t, upgradePlan.Spec.Drain)

	tolerations := []corev1.Toleration{
		{
			Key:      "CriticalAddonsOnly",
			Operator: "Equal",
			Value:    "true",
			Effect:   "NoExecute",
		},
		{
			Key:      ControlPlaneLabel,
			Operator: "Equal",
			Value:    "",
			Effect:   "NoSchedule",
		},
		{
			Key:      "node-role.kubernetes.io/etcd",
			Operator: "Equal",
			Value:    "",
			Effect:   "NoExecute",
		},
	}
	assert.Equal(t, tolerations, upgradePlan.Spec.Tolerations)
}

func TestKubernetesWorkerPlan_RKE2(t *testing.T) {
	version := "v1.30.2+rke2r1"
	addLabels := map[string]string{
		"lifecycle.suse.com/x": "z",
	}

	expectedLabels := map[string]string{
		"lifecycle.suse.com/x": "z",
		"k8s-upgrade":          "worker",
	}

	upgradePlan := KubernetesWorkerPlan(planNameSuffix, version, false, addLabels)
	require.NotNil(t, upgradePlan)

	assert.Equal(t, "Plan", upgradePlan.TypeMeta.Kind)
	assert.Equal(t, "upgrade.cattle.io/v1", upgradePlan.TypeMeta.APIVersion)

	assert.Equal(t, "workers-v1-30-2-rke2r1-abcdef", upgradePlan.ObjectMeta.Name)
	assert.Equal(t, "cattle-system", upgradePlan.ObjectMeta.Namespace)
	assert.Equal(t, expectedLabels, upgradePlan.ObjectMeta.Labels)

	require.Len(t, upgradePlan.Spec.NodeSelector.MatchLabels, 0)
	require.Len(t, upgradePlan.Spec.NodeSelector.MatchExpressions, 1)

	matchExpression := upgradePlan.Spec.NodeSelector.MatchExpressions[0]
	assert.Equal(t, "node-role.kubernetes.io/control-plane", matchExpression.Key)
	assert.EqualValues(t, "NotIn", matchExpression.Operator)
	assert.Equal(t, []string{"true"}, matchExpression.Values)

	prepareContainer := upgradePlan.Spec.Prepare
	require.NotNil(t, prepareContainer)
	assert.Equal(t, "rancher/rke2-upgrade", prepareContainer.Image)
	assert.Empty(t, prepareContainer.Command)
	assert.Equal(t, []string{"prepare", "control-plane-v1-30-2-rke2r1-abcdef"}, prepareContainer.Args)

	upgradeContainer := upgradePlan.Spec.Upgrade
	require.NotNil(t, upgradeContainer)
	assert.Equal(t, "rancher/rke2-upgrade", upgradeContainer.Image)
	assert.Empty(t, upgradeContainer.Command)
	assert.Empty(t, upgradeContainer.Args)

	assert.Equal(t, version, upgradePlan.Spec.Version)
	assert.EqualValues(t, 1, upgradePlan.Spec.Concurrency)
	assert.True(t, upgradePlan.Spec.Cordon)

	assert.Equal(t, "system-upgrade-controller", upgradePlan.Spec.ServiceAccountName)
	assert.Nil(t, upgradePlan.Spec.Drain)

	assert.Len(t, upgradePlan.Spec.Tolerations, 0)
}

func TestKubernetesWorkerPlan_K3s(t *testing.T) {
	version := "v1.30.2+k3s1"
	addLabels := map[string]string{
		"lifecycle.suse.com/x": "z",
	}

	expectedLabels := map[string]string{
		"lifecycle.suse.com/x": "z",
		"k8s-upgrade":          "worker",
	}

	upgradePlan := KubernetesWorkerPlan(planNameSuffix, version, false, addLabels)
	require.NotNil(t, upgradePlan)

	assert.Equal(t, "Plan", upgradePlan.TypeMeta.Kind)
	assert.Equal(t, "upgrade.cattle.io/v1", upgradePlan.TypeMeta.APIVersion)

	assert.Equal(t, "workers-v1-30-2-k3s1-abcdef", upgradePlan.ObjectMeta.Name)
	assert.Equal(t, "cattle-system", upgradePlan.ObjectMeta.Namespace)
	assert.Equal(t, expectedLabels, upgradePlan.ObjectMeta.Labels)

	require.Len(t, upgradePlan.Spec.NodeSelector.MatchLabels, 0)
	require.Len(t, upgradePlan.Spec.NodeSelector.MatchExpressions, 1)

	matchExpression := upgradePlan.Spec.NodeSelector.MatchExpressions[0]
	assert.Equal(t, "node-role.kubernetes.io/control-plane", matchExpression.Key)
	assert.EqualValues(t, "NotIn", matchExpression.Operator)
	assert.Equal(t, []string{"true"}, matchExpression.Values)

	prepareContainer := upgradePlan.Spec.Prepare
	require.NotNil(t, prepareContainer)
	assert.Equal(t, "rancher/k3s-upgrade", prepareContainer.Image)
	assert.Empty(t, prepareContainer.Command)
	assert.Equal(t, []string{"prepare", "control-plane-v1-30-2-k3s1-abcdef"}, prepareContainer.Args)

	upgradeContainer := upgradePlan.Spec.Upgrade
	require.NotNil(t, upgradeContainer)
	assert.Equal(t, "rancher/k3s-upgrade", upgradeContainer.Image)
	assert.Empty(t, upgradeContainer.Command)
	assert.Empty(t, upgradeContainer.Args)

	assert.Equal(t, version, upgradePlan.Spec.Version)
	assert.EqualValues(t, 1, upgradePlan.Spec.Concurrency)
	assert.True(t, upgradePlan.Spec.Cordon)

	assert.Equal(t, "system-upgrade-controller", upgradePlan.Spec.ServiceAccountName)
	assert.Nil(t, upgradePlan.Spec.Drain)

	assert.Len(t, upgradePlan.Spec.Tolerations, 0)
}
