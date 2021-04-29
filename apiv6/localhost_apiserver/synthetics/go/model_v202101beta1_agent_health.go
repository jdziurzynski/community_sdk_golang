/*
 * Synthetics Monitoring API
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 202101beta1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package syntheticsstub

type V202101beta1AgentHealth struct {
	Agent V202101beta1Agent `json:"agent,omitempty"`

	Health []V202101beta1HealthMoment `json:"health,omitempty"`

	OverallHealth V202101beta1Health `json:"overallHealth,omitempty"`
}
