// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package v2alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	iputil "github.com/cilium/cilium/pkg/ip"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:categories={cilium},singular="ciliumsubnettopology",path="ciliumsubnettopologies",scope="Cluster",shortName={cst}
// +kubebuilder:object:root=true
// +kubebuilder:storageversion

// CiliumSubnetTopology describes the cluster-wide subnet topology used with
// hybrid routing.
type CiliumSubnetTopology struct {
	// +deepequal-gen=false
	metav1.TypeMeta `json:",inline"`
	// +deepequal-gen=false
	// +kubebuilder:validation:Optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +kubebuilder:validation:Required
	Spec CiliumSubnetTopologySpec `json:"spec"`
}

type CiliumSubnetTopologySpec struct {
	// SubnetGroups is the list of named subnet groups. Packets are natively
	// routed within the same group.
	//
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	// +listType=map
	// +listMapKey=name
	SubnetGroups []SubnetGroup `json:"subnetGroups"`
}

type SubnetGroup struct {
	// Name is the unique, human-readable name of the subnet group.
	//
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// CIDRs is the list of CIDRs that belong to this subnet group.
	//
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	CIDRs []iputil.Prefix `json:"cidrs"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +deepequal-gen=false
type CiliumSubnetTopologyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is a list of CiliumSubnetTopologies.
	Items []CiliumSubnetTopology `json:"items"`
}
