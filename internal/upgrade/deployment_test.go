package upgrade

import (
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func TestIsDeploymentReady(t *testing.T) {
	d := &appsv1.Deployment{
		Status: appsv1.DeploymentStatus{
			ReadyReplicas: 1,
			Replicas:      1,
		},
	}

	assert.True(t, IsDeploymentReady(d))
	d.Status.ReadyReplicas = 0
	assert.False(t, IsDeploymentReady(d))
}

func TestContainsContainerImagesStrict(t *testing.T) {
	additionalStrictTest := ContainersImageTest{
		Name: "Container image registry mismatch",
		Containers: []corev1.Container{
			{
				Name:  "coredns",
				Image: "foo.bar:8080/rancher/mirrored-coredns-coredns:1.11.3",
			},
		},
		ContainerImages: map[string]string{
			"coredns": "rancher/mirrored-coredns-coredns:1.11.3",
		},
		ExpectedResult: false,
	}

	allTests := append(getDefaultContainerImageTests(), additionalStrictTest)

	for _, test := range allTests {
		t.Run(test.Name, func(t *testing.T) {
			if test.ExpectedResult {
				assert.True(t, ContainsContainerImages(test.Containers, test.ContainerImages, true))
			} else {
				assert.False(t, ContainsContainerImages(test.Containers, test.ContainerImages, true))
			}
		})
	}
}

func TestContainsContainerImagesLenient(t *testing.T) {
	additionalLenientTest := ContainersImageTest{
		Name: "Different container image registry",
		Containers: []corev1.Container{
			{
				Name:  "coredns",
				Image: "foo.bar:8080/rancher/mirrored-coredns-coredns:1.11.3",
			},
		},
		ContainerImages: map[string]string{
			"coredns": "rancher/mirrored-coredns-coredns:1.11.3",
		},
		ExpectedResult: true,
	}

	allTests := append(getDefaultContainerImageTests(), additionalLenientTest)

	for _, test := range allTests {
		t.Run(test.Name, func(t *testing.T) {
			if test.ExpectedResult {
				assert.True(t, ContainsContainerImages(test.Containers, test.ContainerImages, false))
			} else {
				assert.False(t, ContainsContainerImages(test.Containers, test.ContainerImages, false))
			}
		})
	}
}

type ContainersImageTest struct {
	Name            string
	Containers      []corev1.Container
	ContainerImages map[string]string
	ExpectedResult  bool
}

func getDefaultContainerImageTests() []ContainersImageTest {
	return []ContainersImageTest{
		{
			Name: "Missing core container",
			Containers: []corev1.Container{
				{
					Name:  "foo",
					Image: "bar/baz:0.0.0",
				},
			},
			ContainerImages: map[string]string{
				"coredns": "rancher/mirrored-coredns-coredns:1.11.3",
			},
			ExpectedResult: false,
		},
		{
			Name: "Container image version mismatch",
			Containers: []corev1.Container{
				{
					Name:  "coredns",
					Image: "rancher/mirrored-coredns-coredns:1.11.2",
				},
			},
			ContainerImages: map[string]string{
				"coredns": "rancher/mirrored-coredns-coredns:1.11.3",
			},
			ExpectedResult: false,
		},
		{
			Name: "Matching container image without sidecar injection",
			Containers: []corev1.Container{
				{
					Name:  "coredns",
					Image: "rancher/mirrored-coredns-coredns:1.11.3",
				},
			},
			ContainerImages: map[string]string{
				"coredns": "rancher/mirrored-coredns-coredns:1.11.3",
			},
			ExpectedResult: true,
		},
		{
			Name: "Matching container image with sidecar injection",
			Containers: []corev1.Container{
				{
					Name:  "coredns",
					Image: "rancher/mirrored-coredns-coredns:1.11.3",
				},
				{
					Name:  "sidecar",
					Image: "foo/bar:0.0.0",
				},
			},
			ContainerImages: map[string]string{
				"coredns": "rancher/mirrored-coredns-coredns:1.11.3",
			},
			ExpectedResult: true,
		},
	}
}
