package service

import (
	"context"
	"errors"
	openapi "github.com/ghkadim/highload_architect/generated/go-server/go"
	"net/http"
)

// FriendDeleteUserIdPut -
func (s *ApiService) FriendDeleteUserIdPut(ctx context.Context, userId string) (openapi.ImplResponse, error) {
	// TODO - update FriendDeleteUserIdPut with the required logic for this service method.
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

	return openapi.Response(http.StatusNotImplemented, nil), errors.New("FriendDeleteUserIdPut method not implemented")
}

// FriendSetUserIdPut -
func (s *ApiService) FriendSetUserIdPut(ctx context.Context, userId string) (openapi.ImplResponse, error) {
	// TODO - update FriendSetUserIdPut with the required logic for this service method.
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

	return openapi.Response(http.StatusNotImplemented, nil), errors.New("FriendSetUserIdPut method not implemented")
}
