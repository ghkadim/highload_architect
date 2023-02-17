package service

import (
	"context"
	"errors"
	openapi "github.com/ghkadim/highload_architect/generated/go-server/go"
	"net/http"
)

// PostCreatePost -
func (s *ApiService) PostCreatePost(ctx context.Context, postCreatePostRequest openapi.PostCreatePostRequest) (openapi.ImplResponse, error) {
	// TODO - update PostCreatePost with the required logic for this service method.
	// Add api_default_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, string{}) or use other options such as http.Ok ...
	//return Response(200, string{}), nil

	//TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	//return Response(400, nil),nil

	//TODO: Uncomment the next line to return response Response(401, {}) or use other options such as http.Ok ...
	//return Response(401, nil),nil

	//TODO: Uncomment the next line to return response Response(500, LoginPost500Response{}) or use other options such as http.Ok ...
	//return Response(500, LoginPost500Response{}), nil

	//TODO: Uncomment the next line to return response Response(503, LoginPost500Response{}) or use other options such as http.Ok ...
	//return Response(503, LoginPost500Response{}), nil

	return openapi.Response(http.StatusNotImplemented, nil), errors.New("PostCreatePost method not implemented")
}

// PostDeleteIdPut -
func (s *ApiService) PostDeleteIdPut(ctx context.Context, id string) (openapi.ImplResponse, error) {
	// TODO - update PostDeleteIdPut with the required logic for this service method.
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

	return openapi.Response(http.StatusNotImplemented, nil), errors.New("PostDeleteIdPut method not implemented")
}

// PostFeedGet -
func (s *ApiService) PostFeedGet(ctx context.Context, offset int32, limit int32) (openapi.ImplResponse, error) {
	// TODO - update PostFeedGet with the required logic for this service method.
	// Add api_default_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []Post{}) or use other options such as http.Ok ...
	//return Response(200, []Post{}), nil

	//TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	//return Response(400, nil),nil

	//TODO: Uncomment the next line to return response Response(401, {}) or use other options such as http.Ok ...
	//return Response(401, nil),nil

	//TODO: Uncomment the next line to return response Response(500, LoginPost500Response{}) or use other options such as http.Ok ...
	//return Response(500, LoginPost500Response{}), nil

	//TODO: Uncomment the next line to return response Response(503, LoginPost500Response{}) or use other options such as http.Ok ...
	//return Response(503, LoginPost500Response{}), nil

	return openapi.Response(http.StatusNotImplemented, nil), errors.New("PostFeedGet method not implemented")
}

// PostGetIdGet -
func (s *ApiService) PostGetIdGet(ctx context.Context, id string) (openapi.ImplResponse, error) {
	// TODO - update PostGetIdGet with the required logic for this service method.
	// Add api_default_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, Post{}) or use other options such as http.Ok ...
	//return Response(200, Post{}), nil

	//TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	//return Response(400, nil),nil

	//TODO: Uncomment the next line to return response Response(401, {}) or use other options such as http.Ok ...
	//return Response(401, nil),nil

	//TODO: Uncomment the next line to return response Response(500, LoginPost500Response{}) or use other options such as http.Ok ...
	//return Response(500, LoginPost500Response{}), nil

	//TODO: Uncomment the next line to return response Response(503, LoginPost500Response{}) or use other options such as http.Ok ...
	//return Response(503, LoginPost500Response{}), nil

	return openapi.Response(http.StatusNotImplemented, nil), errors.New("PostGetIdGet method not implemented")
}

// PostUpdatePut -
func (s *ApiService) PostUpdatePut(ctx context.Context, postUpdatePutRequest openapi.PostUpdatePutRequest) (openapi.ImplResponse, error) {
	// TODO - update PostUpdatePut with the required logic for this service method.
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

	return openapi.Response(http.StatusNotImplemented, nil), errors.New("PostUpdatePut method not implemented")
}
