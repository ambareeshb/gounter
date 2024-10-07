package route

import (
	"gounter/api/handler"
	"net/http"
)

// InitRoutes initializes the HTTP routes
func InitRoutes(handler *handler.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	// Define routes for create, update, and delete
	mux.HandleFunc("/counter/create", handler.CreateCounter)
	mux.HandleFunc("/counter/increment", handler.IncrementCounter)
	mux.HandleFunc("/counter/delete", handler.DeleteCounter)

	return mux
}
