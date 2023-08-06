package server

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	Authorize   bool
	HandlerFunc http.HandlerFunc
}
