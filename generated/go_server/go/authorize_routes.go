package openapi

var AuthorizeRoutes = []struct {
	Path   string
	Method string
}{

	{
		Path:   "/post/create",
		Method: "POST",
	},

	{
		Path:   "/dialog/{user_id}/send",
		Method: "POST",
	},

	{
		Path:   "/post/delete/{id}",
		Method: "PUT",
	},

	{
		Path:   "/post/feed",
		Method: "GET",
	},

	{
		Path:   "/friend/set/{user_id}",
		Method: "PUT",
	},

	{
		Path:   "/post/update",
		Method: "PUT",
	},

	{
		Path:   "/friend/delete/{user_id}",
		Method: "PUT",
	},

	{
		Path:   "/dialog/{user_id}/list",
		Method: "GET",
	},
}
