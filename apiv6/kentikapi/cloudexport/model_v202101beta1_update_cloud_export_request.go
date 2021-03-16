/*
 * Cloud Export Admin API
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 202101
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package cloudexport

import (
	"encoding/json"
)

// V202101beta1UpdateCloudExportRequest struct for V202101beta1UpdateCloudExportRequest
type V202101beta1UpdateCloudExportRequest struct {
	Export *V202101beta1CloudExport `json:"export,omitempty"`
}

// NewV202101beta1UpdateCloudExportRequest instantiates a new V202101beta1UpdateCloudExportRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewV202101beta1UpdateCloudExportRequest() *V202101beta1UpdateCloudExportRequest {
	this := V202101beta1UpdateCloudExportRequest{}
	return &this
}

// NewV202101beta1UpdateCloudExportRequestWithDefaults instantiates a new V202101beta1UpdateCloudExportRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewV202101beta1UpdateCloudExportRequestWithDefaults() *V202101beta1UpdateCloudExportRequest {
	this := V202101beta1UpdateCloudExportRequest{}
	return &this
}

// GetExport returns the Export field value if set, zero value otherwise.
func (o *V202101beta1UpdateCloudExportRequest) GetExport() V202101beta1CloudExport {
	if o == nil || o.Export == nil {
		var ret V202101beta1CloudExport
		return ret
	}
	return *o.Export
}

// GetExportOk returns a tuple with the Export field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V202101beta1UpdateCloudExportRequest) GetExportOk() (*V202101beta1CloudExport, bool) {
	if o == nil || o.Export == nil {
		return nil, false
	}
	return o.Export, true
}

// HasExport returns a boolean if a field has been set.
func (o *V202101beta1UpdateCloudExportRequest) HasExport() bool {
	if o != nil && o.Export != nil {
		return true
	}

	return false
}

// SetExport gets a reference to the given V202101beta1CloudExport and assigns it to the Export field.
func (o *V202101beta1UpdateCloudExportRequest) SetExport(v V202101beta1CloudExport) {
	o.Export = &v
}

func (o V202101beta1UpdateCloudExportRequest) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Export != nil {
		toSerialize["export"] = o.Export
	}
	return json.Marshal(toSerialize)
}

type NullableV202101beta1UpdateCloudExportRequest struct {
	value *V202101beta1UpdateCloudExportRequest
	isSet bool
}

func (v NullableV202101beta1UpdateCloudExportRequest) Get() *V202101beta1UpdateCloudExportRequest {
	return v.value
}

func (v *NullableV202101beta1UpdateCloudExportRequest) Set(val *V202101beta1UpdateCloudExportRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableV202101beta1UpdateCloudExportRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableV202101beta1UpdateCloudExportRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableV202101beta1UpdateCloudExportRequest(val *V202101beta1UpdateCloudExportRequest) *NullableV202101beta1UpdateCloudExportRequest {
	return &NullableV202101beta1UpdateCloudExportRequest{value: val, isSet: true}
}

func (v NullableV202101beta1UpdateCloudExportRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableV202101beta1UpdateCloudExportRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
