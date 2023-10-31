package openapi

var AuthorizeRoutes = []struct {
	Path   string
	Method string
}{

	{
		Path:   "/counter/{counter_id}",
		Method: "GET",
	},
}
