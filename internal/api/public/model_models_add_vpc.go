/*
Nexodus API

This is the Nexodus API Server.

API version: 1.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package public

// ModelsAddVPC struct for ModelsAddVPC
type ModelsAddVPC struct {
	Description    string `json:"description,omitempty"`
	Ipv4Cidr       string `json:"ipv4_cidr,omitempty"`
	Ipv6Cidr       string `json:"ipv6_cidr,omitempty"`
	OrganizationId string `json:"organization_id,omitempty"`
	PrivateCidr    bool   `json:"private_cidr,omitempty"`
}
