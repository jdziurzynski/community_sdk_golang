/*
 * Synthetics Monitoring API
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 202101beta1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package syntheticsstub

import (
	"context"
	"fmt"
	"net/http"
)

// SyntheticsAdminServiceApiService is a service that implements the logic for the SyntheticsAdminServiceApiServicer
// This service should implement the business logic for every endpoint for the SyntheticsAdminServiceAPI API.
// Include any external packages or services that will be required by this service.
type SyntheticsAdminServiceApiService struct {
	repo *SyntheticsRepo
}

// NewSyntheticsAdminServiceApiService creates a default api service
func NewSyntheticsAdminServiceApiService(repo *SyntheticsRepo) SyntheticsAdminServiceApiServicer {
	return &SyntheticsAdminServiceApiService{
		repo: repo,
	}
}

// AgentDelete - Delete an agent.
func (s *SyntheticsAdminServiceApiService) AgentDelete(_ context.Context, agentId string) (ImplResponse, error) {
	if err := s.repo.DeleteAgent(agentId); err != nil {
		return errorResponse(http.StatusNotFound, "agent DELETE failed", err), nil
	}

	return Response(http.StatusOK, &map[string]interface{}{}), nil
}

// AgentGet - Get information about an agent.
func (s *SyntheticsAdminServiceApiService) AgentGet(_ context.Context, agentId string) (ImplResponse, error) {
	agent := s.repo.GetAgent(agentId)
	if agent == nil {
		return errorResponse(
			http.StatusNotFound,
			"agent GET failed",
			fmt.Errorf("no such agent %q", agentId),
		), nil
	}

	return Response(http.StatusOK, &V202101beta1GetAgentResponse{Agent: *agent}), nil
}

// AgentPatch - Patch an agent.
func (s *SyntheticsAdminServiceApiService) AgentPatch(_ context.Context, agentId string, body V202101beta1PatchAgentRequest) (ImplResponse, error) {
	body.Agent.Id = agentId

	agent, err := s.repo.PatchAgent(body.Agent)
	if err != nil {
		return errorResponse(http.StatusBadRequest, "agent PATCH failed", err), nil
	}

	return Response(http.StatusOK, &V202101beta1PatchAgentResponse{Agent: *agent}), nil
}

// AgentsList - List Agents.
func (s *SyntheticsAdminServiceApiService) AgentsList(_ context.Context) (ImplResponse, error) {
	return Response(http.StatusOK, &V202101beta1ListAgentsResponse{
		Agents:             s.repo.ListAgents(),
		InvalidAgentsCount: 0,
	}), nil
}

// TestCreate - Create Synthetics Test.
func (s *SyntheticsAdminServiceApiService) TestCreate(_ context.Context, body V202101beta1CreateTestRequest) (ImplResponse, error) {
	test, err := s.repo.CreateTest(body.Test)
	if err != nil {
		return errorResponse(http.StatusBadRequest, "test CREATE failed", err), nil
	}

	return Response(http.StatusOK, &V202101beta1CreateTestResponse{Test: *test}), nil

}

// TestDelete - Delete an Synthetics Test.
func (s *SyntheticsAdminServiceApiService) TestDelete(_ context.Context, id string) (ImplResponse, error) {
	if err := s.repo.DeleteTest(id); err != nil {
		return errorResponse(http.StatusNotFound, "test DELETE failed", err), nil
	}

	return Response(http.StatusOK, &map[string]interface{}{}), nil
}

// TestGet - Get information about Synthetics Test.
func (s *SyntheticsAdminServiceApiService) TestGet(_ context.Context, id string) (ImplResponse, error) {
	test := s.repo.GetTest(id)
	if test == nil {
		return errorResponse(
			http.StatusNotFound,
			"test GET failed",
			fmt.Errorf("no such test %q", id),
		), nil
	}

	return Response(http.StatusOK, &V202101beta1GetTestResponse{Test: *test}), nil
}

// TestPatch - Patch a Synthetics Test.
func (s *SyntheticsAdminServiceApiService) TestPatch(_ context.Context, id string, body V202101beta1PatchTestRequest) (ImplResponse, error) {
	body.Test.Id = id

	test, err := s.repo.PatchTest(body.Test)
	if err != nil {
		return errorResponse(http.StatusBadRequest, "test PATCH failed", err), nil
	}

	return Response(http.StatusOK, &V202101beta1PatchTestResponse{Test: *test}), nil
}

// TestStatusUpdate - Update a test status.
func (s *SyntheticsAdminServiceApiService) TestStatusUpdate(_ context.Context, id string, body V202101beta1SetTestStatusRequest) (ImplResponse, error) {
	if err := s.repo.UpdateTestStatus(id, body.Status); err != nil {
		return errorResponse(http.StatusNotFound, "test status update failed", err), nil
	}

	return Response(http.StatusOK, &map[string]interface{}{}), nil
}

// TestsList - List Synthetics Tests.
func (s *SyntheticsAdminServiceApiService) TestsList(_ context.Context, preset bool) (ImplResponse, error) {
	return Response(
		http.StatusOK,
		&V202101beta1ListTestsResponse{
			Tests:             s.repo.ListTests(),
			InvalidTestsCount: 0,
		},
	), nil
}

func errorResponse(httpCode int, message string, err error) ImplResponse {
	const rpcCodeUnknown = 2 // translation httpCode -> rpcCode not relevant here

	rpcResponse := RpcStatus{
		Code:    rpcCodeUnknown,
		Message: fmt.Sprintf("%v: %v", message, err.Error()),
		Details: []ProtobufAny{},
	}
	return Response(httpCode, rpcResponse)
}
