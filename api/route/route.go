package route

import (
	"gounter/api/auth"
	"gounter/api/handler"
	"net/http"
)

// InitRoutes initializes the HTTP routes
func InitRoutes(handler *handler.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	// Define routes for create, update, and delete
	mux.Handle("/counter/create", auth.AuthorizationMiddleware(http.HandlerFunc(handler.CreateCounter)))
	mux.Handle("/counter/increment", auth.AuthorizationMiddleware(http.HandlerFunc(handler.IncrementCounter)))
	mux.Handle("/counter/delete", auth.AuthorizationMiddleware(http.HandlerFunc(handler.DeleteCounter)))

	return mux
}
