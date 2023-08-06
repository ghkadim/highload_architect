package openapi

import (
	"context"
	"errors"

	openapi "github.com/ghkadim/highload_architect/generated/app/go_server/go"
	"github.com/ghkadim/highload_architect/internal/logger"
	"github.com/ghkadim/highload_architect/internal/models"
)

type service interface {
	UserRegister(ctx context.Context, user models.User, password string) (models.UserID, error)
	UserGet(ctx context.Context, id models.UserID) (models.User, error)
	UserSearch(ctx context.Context, firstName, secondName string) ([]models.User, error)
	UserLogin(ctx context.Context, id models.UserID, password string) (string, error)
	FriendAdd(ctx context.Context, userID1, userID2 models.UserID) error
	FriendDelete(ctx context.Context, userID1, userID2 models.UserID) error
	PostAdd(ctx context.Context, text string, author models.UserID) (models.PostID, error)
	PostUpdate(ctx context.Context, userID models.UserID, postID models.PostID, text string) error
	PostDelete(ctx context.Context, userID models.UserID, postID models.PostID) error
	PostGet(ctx context.Context, postID models.PostID) (models.Post, error)
	PostFeed(ctx context.Context, userID models.UserID, offset, limit int) ([]models.Post, error)
	DialogSend(ctx context.Context, message models.DialogMessage) error
	DialogList(ctx context.Context, userID1, userID2 models.UserID) ([]models.DialogMessage, error)
}

type session interface {
	ParseToken(ctx context.Context, tokenStr string) (models.UserID, error)
}

type apiController struct {
	service service
	session session
}

func NewController(
	service service,
	session session,
) *apiController {
	return &apiController{
		service: service,
		session: session,
	}
}

// UserGetIdGet -
func (c *apiController) UserGetIdGet(ctx context.Context, id string) (openapi.ImplResponse, error) {
	user, err := c.service.UserGet(ctx, models.UserID(id))
	if err != nil {
		return errorResponse(err)
	}

	return successResponse(userModelToAPI(user))
}

// UserRegisterPost -
func (c *apiController) UserRegisterPost(ctx context.Context, user openapi.UserRegisterPostRequest) (openapi.ImplResponse, error) {
	id, err := c.service.UserRegister(ctx, models.User{
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		Age:        &user.Age,
		Biography:  &user.Biography,
		City:       &user.City,
	}, user.Password)
	if err != nil {
		return errorResponse(err)
	}

	return successResponse(openapi.UserRegisterPost200Response{UserId: string(id)})
}

// UserSearchGet -
func (c *apiController) UserSearchGet(ctx context.Context, firstName string, lastName string) (openapi.ImplResponse, error) {
	if firstName == "" && lastName == "" {
		return openapi.Response(400, "last_name or first_name should not be empty"), nil
	}

	users, err := c.service.UserSearch(ctx, firstName, lastName)
	if err != nil {
		return errorResponse(err)
	}

	apiUsers := make([]openapi.User, 0, len(users))
	for i := range users {
		apiUsers = append(apiUsers, userModelToAPI(users[i]))
	}

	return successResponse(apiUsers)
}

func (c *apiController) LoginPost(ctx context.Context, loginPostRequest openapi.LoginPostRequest) (openapi.ImplResponse, error) {
	token, err := c.service.UserLogin(ctx, models.UserID(loginPostRequest.Id), loginPostRequest.Password)
	if err != nil {
		return errorResponse(err)
	}
	return successResponse(openapi.LoginPost200Response{Token: token})
}

// PostCreatePost -
func (c *apiController) PostCreatePost(ctx context.Context, postCreatePostRequest openapi.PostCreatePostRequest) (openapi.ImplResponse, error) {
	token := bearerToken(ctx)
	userID, err := c.session.ParseToken(ctx, token)
	if err != nil {
		return errorResponse(err)
	}

	id, err := c.service.PostAdd(ctx, postCreatePostRequest.Text, userID)
	if err != nil {
		return errorResponse(err)
	}

	return successResponse(id)
}

// PostDeleteIdPut -
func (c *apiController) PostDeleteIdPut(ctx context.Context, id string) (openapi.ImplResponse, error) {
	token := bearerToken(ctx)
	userID, err := c.session.ParseToken(ctx, token)
	if err != nil {
		return errorResponse(err)
	}

	err = c.service.PostDelete(ctx, userID, models.PostID(id))
	if err != nil {
		if errors.Is(err, models.ErrPostNotFound) {
			logger.Error("Post already deleted: %v", err)
			return successResponse(nil)
		}
		return errorResponse(err)
	}

	return successResponse(nil)
}

// PostFeedGet -
func (c *apiController) PostFeedGet(ctx context.Context, offset, limit float32) (openapi.ImplResponse, error) {
	token := bearerToken(ctx)
	userID, err := c.session.ParseToken(ctx, token)
	if err != nil {
		return errorResponse(err)
	}

	posts, err := c.service.PostFeed(ctx, userID, int(offset), int(limit))
	if err != nil {
		return errorResponse(err)
	}

	postsResp := make([]openapi.Post, 0, len(posts))
	for _, post := range posts {
		postsResp = append(postsResp, openapi.Post{
			Id:           string(post.ID),
			Text:         post.Text,
			AuthorUserId: string(post.AuthorID),
		})
	}

	return successResponse(postsResp)
}

// PostGetIdGet -
func (c *apiController) PostGetIdGet(ctx context.Context, id string) (openapi.ImplResponse, error) {
	post, err := c.service.PostGet(ctx, models.PostID(id))
	if err != nil {
		return errorResponse(err)
	}

	return successResponse(openapi.Post{
		Id:           string(post.ID),
		Text:         post.Text,
		AuthorUserId: string(post.AuthorID),
	})
}

// PostUpdatePut -
func (c *apiController) PostUpdatePut(ctx context.Context, postUpdatePutRequest openapi.PostUpdatePutRequest) (openapi.ImplResponse, error) {
	token := bearerToken(ctx)
	userID, err := c.session.ParseToken(ctx, token)
	if err != nil {
		return errorResponse(err)
	}

	err = c.service.PostUpdate(ctx, userID, models.PostID(postUpdatePutRequest.Id), postUpdatePutRequest.Text)
	if err != nil {
		return errorResponse(err)
	}

	return successResponse(nil)
}

// FriendDeleteUserIdPut -
func (c *apiController) FriendDeleteUserIdPut(ctx context.Context, friendUserID string) (openapi.ImplResponse, error) {
	token := bearerToken(ctx)
	userID, err := c.session.ParseToken(ctx, token)
	if err != nil {
		return errorResponse(err)
	}

	err = c.service.FriendDelete(ctx, userID, models.UserID(friendUserID))
	if err != nil {
		return errorResponse(err)
	}
	return successResponse(nil)
}

// FriendSetUserIdPut -
func (c *apiController) FriendSetUserIdPut(ctx context.Context, friendUserID string) (openapi.ImplResponse, error) {
	token := bearerToken(ctx)
	userID, err := c.session.ParseToken(ctx, token)
	if err != nil {
		return errorResponse(err)
	}

	err = c.service.FriendAdd(ctx, userID, models.UserID(friendUserID))
	if err != nil {
		return errorResponse(err)
	}
	return successResponse(nil)
}

// DialogUserIdListGet -
func (c *apiController) DialogUserIdListGet(ctx context.Context, userID2 string) (openapi.ImplResponse, error) {
	token := bearerToken(ctx)
	userID1, err := c.session.ParseToken(ctx, token)
	if err != nil {
		return errorResponse(err)
	}

	messages, err := c.service.DialogList(ctx, userID1, models.UserID(userID2))
	if err != nil {
		return errorResponse(err)
	}

	dialogMessages := make([]openapi.DialogMessage, 0, len(messages))
	for _, msg := range messages {
		dialogMessages = append(dialogMessages, openapi.DialogMessage{
			From: string(msg.From),
			To:   string(msg.To),
			Text: msg.Text,
		})
	}
	return successResponse(dialogMessages)
}

// DialogUserIdSendPost -
func (c *apiController) DialogUserIdSendPost(
	ctx context.Context,
	toUserID string,
	dialogUserIdSendPostRequest openapi.DialogUserIdSendPostRequest,
) (openapi.ImplResponse, error) {
	token := bearerToken(ctx)
	fromUserID, err := c.session.ParseToken(ctx, token)
	if err != nil {
		return errorResponse(err)
	}

	err = c.service.DialogSend(
		ctx,
		models.DialogMessage{
			From: fromUserID,
			To:   models.UserID(toUserID),
			Text: dialogUserIdSendPostRequest.Text,
		})
	if err != nil {
		return errorResponse(err)
	}
	return successResponse(nil)
}

func userModelToAPI(user models.User) openapi.User {
	return openapi.User{
		Id:         string(user.ID),
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		Age:        valueOrDefault(user.Age),
		Biography:  valueOrDefault(user.Biography),
		City:       valueOrDefault(user.City),
	}
}

func errorResponse(err error) (openapi.ImplResponse, error) {
	switch {
	case errors.Is(err, models.ErrUserNotFound):
		return openapi.Response(404, nil), err
	case errors.Is(err, models.ErrPostNotFound):
		return openapi.Response(404, nil), err
	default:
		return openapi.Response(500, openapi.LoginPost500Response{}), err
	}
}

func successResponse(body interface{}) (openapi.ImplResponse, error) {
	return openapi.Response(200, body), nil
}

func valueOrDefault[V any](value *V) V {
	if value == nil {
		return *new(V)
	}
	return *value
}

func bearerToken(ctx context.Context) string {
	val := ctx.Value(models.BearerTokenCtxKey)
	if val == nil {
		return ""
	}
	return val.(string)
}
