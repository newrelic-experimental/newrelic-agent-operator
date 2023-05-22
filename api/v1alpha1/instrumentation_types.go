/*
Copyright 2023.

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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// InstrumentationSpec defines the desired state of Instrumentation
type InstrumentationSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	//  Env defines agent specific env vars
	// +optional
	Env []corev1.EnvVar `json:"env,omitempty"`

	// Java defines configuration for java auto-instrumentation.
	// +optional
	Java Java `json:"java,omitempty"`

	// NodeJS defines configuration for nodejs auto-instrumentation.
	// +optional
	NodeJS NodeJS `json:"nodejs,omitempty"`

	// Python defines configuration for python auto-instrumentation.
	// +optional
	Python Python `json:"python,omitempty"`

	// DotNet defines configuration for DotNet auto-instrumentation.
	// +optional
	DotNet DotNet `json:"dotnet,omitempty"`
}

// Java defines Java agent and instrumentation configuration.
type Java struct {
	// Image is a container image with javaagent auto-instrumentation JAR.
	// +optional
	Image string `json:"image,omitempty"`

	// Env defines java specific env vars.
	// If the former var had been defined, then the other vars would be ignored.
	// +optional
	Env []corev1.EnvVar `json:"env,omitempty"`
}

// NodeJS defines NodeJS agent and instrumentation configuration.
type NodeJS struct {
	// Image is a container image with NodeJS agent and auto-instrumentation.
	// +optional
	Image string `json:"image,omitempty"`

	// Env defines nodejs specific env vars.
	// If the former var had been defined, then the other vars would be ignored.
	// +optional
	Env []corev1.EnvVar `json:"env,omitempty"`
}

// Python defines Python agent and instrumentation configuration.
type Python struct {
	// Image is a container image with Python agent and auto-instrumentation.
	// +optional
	Image string `json:"image,omitempty"`

	// Env defines python specific env vars.
	// If the former var had been defined, then the other vars would be ignored.
	// +optional
	Env []corev1.EnvVar `json:"env,omitempty"`
}

type DotNet struct {
	// Image is a container image with DotNet agent and auto-instrumentation.
	// +optional
	Image string `json:"image,omitempty"`

	// Env defines DotNet specific env vars.
	// If the former var had been defined, then the other vars would be ignored.
	// +optional
	Env []corev1.EnvVar `json:"env,omitempty"`
}

// InstrumentationStatus defines the observed state of Instrumentation
type InstrumentationStatus struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=nragent;nragents
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +operator-sdk:csv:customresourcedefinitions:displayName="New Relic Instrumentation"
// +operator-sdk:csv:customresourcedefinitions:resources={{Pod,v1}}

// Instrumentation is the Schema for the instrumentations API
type Instrumentation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InstrumentationSpec   `json:"spec,omitempty"`
	Status InstrumentationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// InstrumentationList contains a list of Instrumentation
type InstrumentationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Instrumentation `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Instrumentation{}, &InstrumentationList{})
}
