package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func (api *API) Server() *http.Server {

	r := mux.NewRouter()

	r.Use(enableCorsMiddleware)
	r.Handle("/", api.IsAuthorized(api.Home))
	r.Handle("/band/{id}", api.IsAuthorized(api.BandPage))
	r.Handle("/delete/{id}", api.IsAuthorizedJWT(api.DeleteUserByID)) // ONLY FOR VALID USER
	r.Handle("/update/{id}", api.IsAuthorizedJWT(api.UpdateUserByID))
	r.Handle("/admin", api.IsAdmin(api.Admin))
	r.HandleFunc("/login", api.LoginPost)
	r.HandleFunc("/signup", api.SingUpPost)
	r.Handle("/posts", api.IsAuthorized(api.Posts)).Methods(http.MethodGet, http.MethodOptions)
	r.Handle("/posts", api.IsAuthorized(api.CreatePost)).Methods(http.MethodPost, http.MethodOptions)
	r.Handle("/posts/{id}", api.IsAuthorized(api.PostByID)).Methods(http.MethodGet, http.MethodOptions)
	r.Handle("/posts/{id}", api.IsAuthorized(api.UpdatePost)).Methods(http.MethodPut, http.MethodOptions)
	r.Handle("/posts/{id}", api.IsAuthorized(api.DeletePost)).Methods(http.MethodDelete, http.MethodOptions)

	server := &http.Server{
		Addr:    ":4000",
		Handler: r,
	}
	fmt.Println("starting server on a :4000")
	return server
}
