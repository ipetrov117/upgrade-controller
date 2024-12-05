package controller

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lifecyclev1alpha1 "github.com/suse-edge/upgrade-controller/api/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestFindMatchingNodes(t *testing.T) {
	nodeLabels := map[string]string{
		"node-x": "z",
	}
	nodeSelector := &metav1.LabelSelector{
		MatchLabels: nodeLabels,
	}

	tests := []struct {
		name          string
		nodeList      *corev1.NodeList
		expectedNodes []string
		expectedErr   string
	}{
		{
			name: "All nodes match",
			nodeList: &corev1.NodeList{
				Items: []corev1.Node{
					{ObjectMeta: metav1.ObjectMeta{Name: "node-1", Labels: nodeLabels}},
					{ObjectMeta: metav1.ObjectMeta{Name: "node-2", Labels: nodeLabels}},
					{ObjectMeta: metav1.ObjectMeta{Name: "node-3", Labels: nodeLabels}},
				},
			},
			expectedNodes: []string{"node-1", "node-2", "node-3"},
		},
		{
			name: "Some nodes match",
			nodeList: &corev1.NodeList{
				Items: []corev1.Node{
					{ObjectMeta: metav1.ObjectMeta{Name: "node-1", Labels: nodeLabels}},
					{ObjectMeta: metav1.ObjectMeta{Name: "node-2"}},
					{ObjectMeta: metav1.ObjectMeta{Name: "node-3", Labels: nodeLabels}},
				},
			},
			expectedNodes: []string{"node-1", "node-3"},
		},
		{
			name: "No nodes match",
			nodeList: &corev1.NodeList{
				Items: []corev1.Node{
					{ObjectMeta: metav1.ObjectMeta{Name: "node-1"}},
					{ObjectMeta: metav1.ObjectMeta{Name: "node-2"}},
					{ObjectMeta: metav1.ObjectMeta{Name: "node-3"}},
				},
			},
			expectedErr: "none of the nodes match label selector: MatchLabels: map[node-x:z], MatchExpressions: []",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nodes, err := findMatchingNodes(test.nodeList, nodeSelector)
			if test.expectedErr != "" {
				require.EqualError(t, err, test.expectedErr)
				assert.Nil(t, nodes)
				return
			}

			require.NoError(t, err)
			require.Len(t, nodes, len(test.expectedNodes))
			for _, expected := range test.expectedNodes {
				assert.True(t, slices.ContainsFunc(nodes, func(actual corev1.Node) bool {
					return actual.Name == expected
				}))
			}
		})
	}
}

func TestIsKubernetesUpgraded(t *testing.T) {
	const kubernetesVersion = "v1.30.3+k3s1"

	tests := []struct {
		name            string
		nodes           []corev1.Node
		expectedUpgrade bool
	}{
		{
			name: "All nodes upgraded",
			nodes: []corev1.Node{
				{
					Status: corev1.NodeStatus{
						Conditions: []corev1.NodeCondition{{Type: corev1.NodeReady, Status: corev1.ConditionTrue}},
						NodeInfo:   corev1.NodeSystemInfo{KubeletVersion: "v1.30.3+k3s1"}},
				},
				{
					Status: corev1.NodeStatus{
						Conditions: []corev1.NodeCondition{{Type: corev1.NodeReady, Status: corev1.ConditionTrue}},
						NodeInfo:   corev1.NodeSystemInfo{KubeletVersion: "v1.30.3+k3s1"}},
				},
				{
					Status: corev1.NodeStatus{
						Conditions: []corev1.NodeCondition{{Type: corev1.NodeReady, Status: corev1.ConditionTrue}},
						NodeInfo:   corev1.NodeSystemInfo{KubeletVersion: "v1.30.3+k3s1"}},
				},
			},
			expectedUpgrade: true,
		},
		{
			name: "Unschedulable node",
			nodes: []corev1.Node{
				{
					Spec: corev1.NodeSpec{Unschedulable: true},
					Status: corev1.NodeStatus{
						Conditions: []corev1.NodeCondition{{Type: corev1.NodeReady, Status: corev1.ConditionTrue}},
						NodeInfo:   corev1.NodeSystemInfo{KubeletVersion: "v1.30.3+k3s1"}},
				},
				{
					Status: corev1.NodeStatus{
						Conditions: []corev1.NodeCondition{{Type: corev1.NodeReady, Status: corev1.ConditionTrue}},
						NodeInfo:   corev1.NodeSystemInfo{KubeletVersion: "v1.30.3+k3s1"}},
				},
				{
					Status: corev1.NodeStatus{
						Conditions: []corev1.NodeCondition{{Type: corev1.NodeReady, Status: corev1.ConditionTrue}},
						NodeInfo:   corev1.NodeSystemInfo{KubeletVersion: "v1.30.3+k3s1"}},
				},
			},
			expectedUpgrade: false,
		},
		{
			name: "Not ready node",
			nodes: []corev1.Node{
				{
					Status: corev1.NodeStatus{
						Conditions: []corev1.NodeCondition{{Type: corev1.NodeReady, Status: corev1.ConditionFalse}},
						NodeInfo:   corev1.NodeSystemInfo{KubeletVersion: "v1.30.3+k3s1"}},
				},
				{
					Status: corev1.NodeStatus{
						Conditions: []corev1.NodeCondition{{Type: corev1.NodeReady, Status: corev1.ConditionTrue}},
						NodeInfo:   corev1.NodeSystemInfo{KubeletVersion: "v1.30.3+k3s1"}},
				},
				{
					Status: corev1.NodeStatus{
						Conditions: []corev1.NodeCondition{{Type: corev1.NodeReady, Status: corev1.ConditionTrue}},
						NodeInfo:   corev1.NodeSystemInfo{KubeletVersion: "v1.30.3+k3s1"}},
				},
			},
			expectedUpgrade: false,
		},
		{
			name: "Node on older Kubernetes version",
			nodes: []corev1.Node{
				{
					Status: corev1.NodeStatus{
						Conditions: []corev1.NodeCondition{{Type: corev1.NodeReady, Status: corev1.ConditionTrue}},
						NodeInfo:   corev1.NodeSystemInfo{KubeletVersion: "v1.28.12+k3s1"}},
				},
				{
					Status: corev1.NodeStatus{
						Conditions: []corev1.NodeCondition{{Type: corev1.NodeReady, Status: corev1.ConditionTrue}},
						NodeInfo:   corev1.NodeSystemInfo{KubeletVersion: "v1.30.3+k3s1"}},
				},
				{
					Status: corev1.NodeStatus{
						Conditions: []corev1.NodeCondition{{Type: corev1.NodeReady, Status: corev1.ConditionTrue}},
						NodeInfo:   corev1.NodeSystemInfo{KubeletVersion: "v1.28.12+k3s1"}},
				},
			},
			expectedUpgrade: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectedUpgrade, isKubernetesUpgraded(test.nodes, kubernetesVersion))
		})
	}
}

func TestControlPlaneOnlyCluster(t *testing.T) {
	assert.True(t, controlPlaneOnlyCluster(&corev1.NodeList{
		Items: []corev1.Node{
			{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"node-role.kubernetes.io/control-plane": "true"}}},
			{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"node-role.kubernetes.io/control-plane": "true"}}},
		},
	}))

	assert.False(t, controlPlaneOnlyCluster(&corev1.NodeList{
		Items: []corev1.Node{
			{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"node-role.kubernetes.io/control-plane": "true"}}},
			{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"node-role.kubernetes.io/control-plane": "false"}}},
		},
	}))

	assert.False(t, controlPlaneOnlyCluster(&corev1.NodeList{
		Items: []corev1.Node{
			{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"node-role.kubernetes.io/control-plane": "true"}}},
			{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{}}},
		},
	}))

	assert.False(t, controlPlaneOnlyCluster(&corev1.NodeList{
		Items: []corev1.Node{
			{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{}}},
		},
	}))
}

func TestTargetKubernetesVersion(t *testing.T) {
	kubernetes := &lifecyclev1alpha1.Kubernetes{
		K3S: lifecyclev1alpha1.KubernetesDistribution{
			Version: "v1.30.3+k3s1",
		},
		RKE2: lifecyclev1alpha1.KubernetesDistribution{
			Version: "v1.30.3+rke2r1",
		},
	}

	tests := []struct {
		name                 string
		nodes                *corev1.NodeList
		expectedDistribution *lifecyclev1alpha1.KubernetesDistribution
		expectedError        string
	}{
		{
			name:          "Empty node list",
			nodes:         &corev1.NodeList{},
			expectedError: "unable to determine current kubernetes distribution due to empty node list",
		},
		{
			name: "Unsupported Kubernetes distribution",
			nodes: &corev1.NodeList{
				Items: []corev1.Node{{Status: corev1.NodeStatus{NodeInfo: corev1.NodeSystemInfo{KubeletVersion: "v1.30.3"}}}},
			},
			expectedError: "unsupported kubernetes distribution detected in version v1.30.3",
		},
		{
			name: "Target k3s distribution",
			nodes: &corev1.NodeList{
				Items: []corev1.Node{{Status: corev1.NodeStatus{NodeInfo: corev1.NodeSystemInfo{KubeletVersion: "v1.28.12+k3s1"}}}},
			},
			expectedDistribution: &kubernetes.K3S,
		},
		{
			name: "Target RKE2 distribution",
			nodes: &corev1.NodeList{
				Items: []corev1.Node{{Status: corev1.NodeStatus{NodeInfo: corev1.NodeSystemInfo{KubeletVersion: "v1.28.12+rke2r1"}}}},
			},
			expectedDistribution: &kubernetes.RKE2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			k8sDistro, err := targetKubernetesDistribution(test.nodes, kubernetes)
			if test.expectedError != "" {
				require.Error(t, err)
				assert.EqualError(t, err, test.expectedError)
				assert.Nil(t, k8sDistro)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expectedDistribution, k8sDistro)
			}
		})
	}
}

func TestIsK8sCoreDeploymentUpgraded(t *testing.T) {
	readyStatus := appsv1.DeploymentStatus{
		Replicas:      1,
		ReadyReplicas: 1,
	}

	tests := []struct {
		name                   string
		upgardeContainerImages map[string]string
		deployment             *appsv1.Deployment
		expectedResult         bool
	}{
		{
			name: "Deployment not yet upgraded",
			upgardeContainerImages: map[string]string{
				"coredns": "rancher/mirrored-coredns-coredns:1.11.3",
			},
			deployment: &appsv1.Deployment{
				Status: readyStatus,
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "coredns",
									Image: "rancher/mirrored-coredns-coredns:1.11.2",
								},
							},
						},
					},
				},
			},
			expectedResult: false,
		},
		{
			name: "Deployment under upgrade",
			deployment: &appsv1.Deployment{
				Status: appsv1.DeploymentStatus{
					Replicas:      1,
					ReadyReplicas: 0,
				},
			},
			expectedResult: false,
		},
		{
			name: "Deployment missing upgrade container",
			upgardeContainerImages: map[string]string{
				"coredns": "rancher/mirrored-coredns-coredns:1.11.3",
			},
			deployment: &appsv1.Deployment{
				Status: readyStatus,
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "foo",
									Image: "foo/bar:0.0.0",
								},
							},
						},
					},
				},
			},
			expectedResult: false,
		},
		{
			name: "Deployment upgraded - public registry",
			upgardeContainerImages: map[string]string{
				"coredns": "rancher/mirrored-coredns-coredns:1.11.3",
			},
			deployment: &appsv1.Deployment{
				Status: readyStatus,
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "coredns",
									Image: "rancher/mirrored-coredns-coredns:1.11.3",
								},
							},
						},
					},
				},
			},
			expectedResult: true,
		},
		{
			name: "Deployment upgraded - private registry",
			upgardeContainerImages: map[string]string{
				"coredns": "rancher/mirrored-coredns-coredns:1.11.3",
			},
			deployment: &appsv1.Deployment{
				Status: readyStatus,
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "coredns",
									Image: "foo.bar:8080/rancher/mirrored-coredns-coredns:1.11.3",
								},
							},
						},
					},
				},
			},
			expectedResult: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.expectedResult {
				assert.True(t, isK8sCoreDeploymentUpgraded(test.deployment, test.upgardeContainerImages))
			} else {
				assert.False(t, isK8sCoreDeploymentUpgraded(test.deployment, test.upgardeContainerImages))
			}
		})
	}
}
