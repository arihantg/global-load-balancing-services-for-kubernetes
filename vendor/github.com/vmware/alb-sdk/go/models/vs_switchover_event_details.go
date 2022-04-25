// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache License 2.0
package models

// This file is auto-generated.

// VsSwitchoverEventDetails vs switchover event details
// swagger:model VsSwitchoverEventDetails
type VsSwitchoverEventDetails struct {

	// Error messages associated with this Event. Field introduced in 21.1.3.
	ErrorMessage *string `json:"error_message,omitempty"`

	// VIP IPv4 address. Field introduced in 21.1.3.
	IP *string `json:"ip,omitempty"`

	// VIP IPv6 address. Field introduced in 21.1.3.
	Ip6 *string `json:"ip6,omitempty"`

	// Status of Event. Field introduced in 21.1.3.
	RPCStatus *int64 `json:"rpc_status,omitempty"`

	// List of ServiceEngine assigned to this Virtual Service. Field introduced in 21.1.3.
	SeAssigned []*VipSeAssigned `json:"se_assigned,omitempty"`

	// Resources requested/assigned to this Virtual Service. Field introduced in 21.1.3.
	SeRequested *VirtualServiceResource `json:"se_requested,omitempty"`

	// Virtual Service UUID. Field introduced in 21.1.3.
	// Required: true
	VsUUID *string `json:"vs_uuid"`
}
