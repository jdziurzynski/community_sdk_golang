/*
 * Cloud Export Admin API
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 202101beta1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package cloudexportstub

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

// CloudExportAdminServiceApiService is a service that implements the logic for the CloudExportAdminServiceApiServicer
// This service should implement the business logic for every endpoint for the CloudExportAdminServiceAPI API.
// Include any external packages or services that will be required by this service.
type CloudExportAdminServiceApiService struct {
	repo *CloudExportRepo
}

// NewCloudExportAdminServiceApiService creates a default api service
func NewCloudExportAdminServiceApiService(repo *CloudExportRepo) CloudExportAdminServiceApiServicer {
	return &CloudExportAdminServiceApiService{
		repo: repo,
	}
}

// ExportCreate - Create Cloud Export.
func (s *CloudExportAdminServiceApiService) ExportCreate(_ context.Context, body V202101beta1CreateCloudExportRequest) (ImplResponse, error) {
	export, err := s.repo.Create(body.Export)
	if err != nil {
		return errorResponse(http.StatusBadRequest, "cloud export CREATE failed", err), nil
	}

	return Response(http.StatusOK, &V202101beta1CreateCloudExportResponse{Export: *export}), nil
}

// ExportDelete - Delete an export.
func (s *CloudExportAdminServiceApiService) ExportDelete(_ context.Context, exportId string) (ImplResponse, error) {
	if err := s.repo.Delete(exportId); err != nil {
		return errorResponse(http.StatusNotFound, "cloud export DELETE failed", err), nil
	}

	return Response(http.StatusOK, &map[string]interface{}{}), nil
}

// ExportGet - Get information about an export.
func (s *CloudExportAdminServiceApiService) ExportGet(_ context.Context, exportId string) (ImplResponse, error) {
	export := s.repo.Get(exportId)
	if export == nil {
		return errorResponse(
			http.StatusNotFound,
			"cloud export GET failed",
			fmt.Errorf("no such cloud export %q", exportId),
		), nil
	}

	return Response(http.StatusOK, &V202101beta1GetCloudExportResponse{Export: *export}), nil
}

// ExportList - List Cloud Export.
func (s *CloudExportAdminServiceApiService) ExportList(_ context.Context) (ImplResponse, error) {
	return Response(http.StatusOK, &V202101beta1ListCloudExportResponse{
		Exports:             s.repo.List(),
		InvalidExportsCount: 0,
	}), nil
}

// ExportPatch - Patch an export.
func (s *CloudExportAdminServiceApiService) ExportPatch(_ context.Context, exportId string, body V202101beta1PatchCloudExportRequest) (ImplResponse, error) {
	// TODO - update ExportPatch with the required logic for this service method.

	// TODO: Uncomment the next line to return response Response(200, V202101beta1PatchCloudExportResponse{}) or use other options such as http.Ok ...
	//return Response(200, V202101beta1PatchCloudExportResponse{}), nil

	// TODO: Uncomment the next line to return response Response(0, GooglerpcStatus{}) or use other options such as http.Ok ...
	//return Response(0, GooglerpcStatus{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("ExportPatch method not implemented")
}

// ExportUpdate - Update an export.
func (s *CloudExportAdminServiceApiService) ExportUpdate(_ context.Context, exportId string, body V202101beta1UpdateCloudExportRequest) (ImplResponse, error) {
	body.Export.Id = exportId

	export, err := s.repo.Update(body.Export)
	if err != nil {
		return errorResponse(http.StatusBadRequest, "cloud export UPDATE failed", err), nil
	}

	return Response(http.StatusOK, &V202101beta1UpdateCloudExportResponse{Export: *export}), nil

}

func errorResponse(httpCode int, message string, err error) ImplResponse {
	const grpcCodeUnknown = 2 // translation httpCode -> grpcCode not relevant here

	grpcResponse := GooglerpcStatus{
		Code:    grpcCodeUnknown,
		Message: fmt.Sprintf("%v: %v", message, err.Error()),
		Details: []ProtobufAny{},
	}
	return Response(httpCode, grpcResponse)
}
