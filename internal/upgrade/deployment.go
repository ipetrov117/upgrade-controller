package upgrade

import (
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func IsDeploymentReady(dep *appsv1.Deployment) bool {
	return dep.Status.Replicas == dep.Status.ReadyReplicas
}

func ContainsContainerImages(containers []corev1.Container, contains map[string]string, strict bool) bool {
	foundContainers := 0
	for _, container := range containers {
		image, ok := contains[container.Name]
		if !ok {
			// Skip containers that are not in the 'contains' map;
			// for use-cases where additional sidecar containers are
			// dynamically added to the resource upon creation.
			continue
		}
		foundContainers++

		if strict && container.Image != image {
			// Strict image comparison.
			return false
		}

		if !strict && !strings.Contains(container.Image, image) {
			// Lenient image comparison;
			// for use-cases where the image registry may change
			// based on the environment use-case (e.g. private regisrty).
			return false
		}
	}

	return foundContainers == len(contains)
}
