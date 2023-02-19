package service

import (
	"context"
	"errors"
	openapi "github.com/ghkadim/highload_architect/generated/go_server/go"
	"net/http"
)

// DialogUserIdListGet -
func (s *ApiService) DialogUserIdListGet(ctx context.Context, userId string) (openapi.ImplResponse, error) {
	// TODO - update DialogUserIdListGet with the required logic for this service method.
	// Add api_default_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []DialogMessage{}) or use other options such as http.Ok ...
	//return Response(200, []DialogMessage{}), nil

	//TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	//return Response(400, nil),nil

	//TODO: Uncomment the next line to return response Response(401, {}) or use other options such as http.Ok ...
	//return Response(401, nil),nil

	//TODO: Uncomment the next line to return response Response(500, LoginPost500Response{}) or use other options such as http.Ok ...
	//return Response(500, LoginPost500Response{}), nil

	//TODO: Uncomment the next line to return response Response(503, LoginPost500Response{}) or use other options such as http.Ok ...
	//return Response(503, LoginPost500Response{}), nil

	return openapi.Response(http.StatusNotImplemented, nil), errors.New("DialogUserIdListGet method not implemented")
}

// DialogUserIdSendPost -
func (s *ApiService) DialogUserIdSendPost(
	ctx context.Context,
	userId string,
	dialogUserIdSendPostRequest openapi.DialogUserIdSendPostRequest,
) (openapi.ImplResponse, error) {
	// TODO - update DialogUserIdSendPost with the required logic for this service method.
	// Add api_default_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	//return Response(400, nil),nil

	//TODO: Uncomment the next line to return response Response(401, {}) or use other options such as http.Ok ...
	//return Response(401, nil),nil

	//TODO: Uncomment the next line to return response Response(500, LoginPost500Response{}) or use other options such as http.Ok ...
	//return Response(500, LoginPost500Response{}), nil

	//TODO: Uncomment the next line to return response Response(503, LoginPost500Response{}) or use other options such as http.Ok ...
	//return Response(503, LoginPost500Response{}), nil

	return openapi.Response(http.StatusNotImplemented, nil), errors.New("DialogUserIdSendPost method not implemented")
}
