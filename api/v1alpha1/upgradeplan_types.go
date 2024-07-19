/*
Copyright 2024.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	KubernetesUpgradedCondition = "KubernetesUpgraded"

	// UpgradePending indicates that the upgrade process has not begun.
	UpgradePending = "Pending"

	// UpgradeInProgress indicates that the upgrade process has started.
	UpgradeInProgress = "InProgress"

	// UpgradeSucceeded indicates that the upgrade process has been successful.
	UpgradeSucceeded = "Succeeded"

	// UpgradeFailed indicates that an error occurred during the upgrade process.
	UpgradeFailed = "Failed"
)

// UpgradePlanSpec defines the desired state of UpgradePlan
type UpgradePlanSpec struct {
	// ReleaseVersion specifies the target version for platform upgrade.
	// The version format is X.Y.Z, for example "3.0.2".
	ReleaseVersion string `json:"releaseVersion"`
}

// UpgradePlanStatus defines the observed state of UpgradePlan
type UpgradePlanStatus struct {
	// +listType=map
	// +listMapKey=type
	// +patchStrategy=merge
	// +patchMergeKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// UpgradePlan is the Schema for the upgradeplans API
type UpgradePlan struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   UpgradePlanSpec   `json:"spec,omitempty"`
	Status UpgradePlanStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// UpgradePlanList contains a list of UpgradePlan
type UpgradePlanList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []UpgradePlan `json:"items"`
}

func init() {
	SchemeBuilder.Register(&UpgradePlan{}, &UpgradePlanList{})
}