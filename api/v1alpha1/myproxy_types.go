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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MyProxySpec defines the desired state of MyProxy
type MyProxySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Name
	Name string `json:"name,omitempty"`
}

// MyProxyStatus defines the observed state of MyProxy
type MyProxyStatus struct {
	// pod名称
	PodNames []string `json:"pod_names"`
	// pod条件
	Conditions []metav1.Condition `json:"conditions"` 
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MyProxy is the Schema for the myproxies API
type MyProxy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MyProxySpec   `json:"spec,omitempty"`
	Status MyProxyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MyProxyList contains a list of MyProxy
type MyProxyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MyProxy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MyProxy{}, &MyProxyList{})
}
