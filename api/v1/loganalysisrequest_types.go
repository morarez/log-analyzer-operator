/*
Copyright 2025.

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

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// LogAnalysisRequestSpec defines the desired state
type LogAnalysisRequestSpec struct {
	ObjectRef corev1.ObjectReference `json:"objectRef"`
	TailLines *int64                 `json:"tailLines,omitempty"`
}

// LogAnalysisRequestStatus defines the observed state
type LogAnalysisRequestStatus struct {
	Diagnosis string      `json:"diagnosis,omitempty"`
	Resolved  bool        `json:"resolved,omitempty"`
	Timestamp metav1.Time `json:"timestamp,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// LogAnalysisRequest is the Schema for the loganalysisrequests API
type LogAnalysisRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LogAnalysisRequestSpec   `json:"spec,omitempty"`
	Status LogAnalysisRequestStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// LogAnalysisRequestList contains a list of LogAnalysisRequest
type LogAnalysisRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LogAnalysisRequest `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LogAnalysisRequest{}, &LogAnalysisRequestList{})
}
