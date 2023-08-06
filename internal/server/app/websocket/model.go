package websocket

type post struct {
	Id           string `json:"postId"`
	Text         string `json:"postText"`
	AuthorUserId string `json:"author_user_id"`
}
