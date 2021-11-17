/*
Copyright 2021.

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
	"github.com/kubesphere-sigs/ks-releaser-operator/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ReleaserControllerSpec defines the desired state of ReleaserController
type ReleaserControllerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Image is an example field of ReleaserController. Edit releasercontroller_types.go to remove/update
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Image",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:fieldGroup:Core","urn:alm:descriptor:com.tectonic.ui:text"}
	Image string `json:"image,omitempty"`

	// Version is the Dex container image tag.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Version",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:fieldGroup:Core","urn:alm:descriptor:com.tectonic.ui:text"}
	Version string `json:"version,omitempty"`

	// Webhook is the Dex container image tag.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Webhook",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:fieldGroup:Core","urn:alm:descriptor:com.tectonic.ui:text"}
	Webhook bool `json:"webhook,omitempty"`
}

// ReleaserControllerStatus defines the observed state of ReleaserController
type ReleaserControllerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ReleaserController is the Schema for the releasercontrollers API
type ReleaserController struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ReleaserControllerSpec   `json:"spec,omitempty"`
	Status ReleaserControllerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ReleaserControllerList contains a list of ReleaserController
type ReleaserControllerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ReleaserController `json:"items"`
}

// IsDeletionFinalizerPresent checks if the instance has deletion finalizer
func (r *ReleaserController) IsDeletionFinalizerPresent() bool {
	for _, finalizer := range r.GetFinalizers() {
		if finalizer == common.DeletionFinalizer {
			return true
		}
	}
	return false
}

func init() {
	SchemeBuilder.Register(&ReleaserController{}, &ReleaserControllerList{})
}
