/*
 * Synthetics Monitoring API
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 202101beta1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package syntheticsstub

type V202101beta1SetTestStatusRequest struct {
	Id string `json:"id,omitempty"`

	Status V202101beta1TestStatus `json:"status,omitempty"`
}
