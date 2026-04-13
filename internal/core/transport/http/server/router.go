package core_http_server

import (
	"fmt"
	"net/http"
)

type ApiVersion string

var (
	ApiVersion1 = ApiVersion("1")
	ApiVersion2 = ApiVersion("2")
	ApiVersion3 = ApiVersion("3")
)

type ApiVersionRouter struct {
	*http.ServeMux
	apiVersion ApiVersion
}

func NewApiVersionRouter(
	apiVersion ApiVersion,
) *ApiVersionRouter {
	return &ApiVersionRouter{
		ServeMux: http.NewServeMux(),
		apiVersion: apiVersion,
	}
}

func (r *ApiVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)

		r.Handle(pattern, route.Handler)
	}
}