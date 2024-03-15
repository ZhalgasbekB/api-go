package main

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/internal/models"
	"net/http"
	"time"
)

func (api *API) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		api.Error(w, http.StatusNotFound, "Not found.")
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		api.Error(w, http.StatusMethodNotAllowed, "Method not allowed.")
		return
	}
	WriteJSON(w, http.StatusOK, api.Artists)
}

func (api *API) BandPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		api.Error(w, http.StatusMethodNotAllowed, "Method not allowed.")
		return
	}
	id, err := getID(r)
	if err != nil {
		api.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	ex := false
	for _, v := range api.Artists {
		if v.ID == id {
			ex = true
			WriteJSON(w, http.StatusOK, v)
		}
	}
	if !ex {
		api.Error(w, http.StatusNotFound, "Artist Doesn't exist")
		return
	}
}

func (api *API) LoginPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		api.Error(w, http.StatusMethodNotAllowed, "Method not allowed.")
		return
	}
	var login models.User
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		api.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := api.DB.UserByEmail(login.Email, login.HashPassword)
	if err != nil {
		api.Error(w, http.StatusNotFound, err.Error())
		return
	}
	api.CacheUser = user
	tokenString, err := api.GenerateJWT(user.Email, user.IsAdmin, user.Name, user.ID)
	if err != nil {
		api.Error(w, http.StatusNotFound, fmt.Errorf("Error generating token: %v", err).Error())
		return
	}
	posts, err := api.DB.Posts(user.ID)
	if err != nil {
		api.Error(w, http.StatusNotFound, err.Error())
		return
	}

	api.UserPosts = &models.UsersPosts{user, MapOfPosts(posts)}
	res := &models.JWToken{Token: tokenString}
	WriteJSON(w, http.StatusOK, res)
}
func (api *API) SingUpPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		api.Error(w, http.StatusMethodNotAllowed, "Method not allowed.")
		return
	}
	sign := models.User{}
	if err := json.NewDecoder(r.Body).Decode(&sign); err != nil {
		api.Error(w, http.StatusMethodNotAllowed, err.Error())
		return
	}
	user, err := models.UserConstructor(sign.Name, sign.Email, sign.HashPassword)
	if err != nil {
		api.Error(w, http.StatusMethodNotAllowed, err.Error())
		return
	}
	id, err := api.DB.CreateUser(user)
	if err != nil {
		api.Error(w, http.StatusMethodNotAllowed, err.Error())
		return
	}
	tokenString, err := api.GenerateJWT(user.Email, user.IsAdmin, user.Name, id)
	if err != nil {
		api.Error(w, http.StatusMethodNotAllowed, fmt.Errorf("Error generating token: %v", err).Error())
		return
	}
	res := &models.JWToken{tokenString}
	WriteJSON(w, http.StatusOK, res)
}

func (api *API) DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.Header().Set("Allow", http.MethodDelete)
		api.Error(w, http.StatusMethodNotAllowed, "Method not allowed.")
		return
	}
	id, err := getID(r)
	if err != nil {
		api.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	user, err := api.DB.DeleteUser(id)
	if err != nil {
		api.Error(w, http.StatusInternalServerError, err.Error())
	}
	WriteJSON(w, http.StatusOK, user)
}

func (api *API) UpdateUserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.Header().Set("Allow", http.MethodPut)
		api.Error(w, http.StatusMethodNotAllowed, "Method not allowed.")
		return
	}

	id, err := getID(r)
	if err != nil {
		api.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := api.DB.UserByID(id)
	if err != nil {
		api.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	var userUpdated models.User
	if err := json.NewDecoder(r.Body).Decode(&userUpdated); err != nil {
		api.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if IsAdmin {
		userUpdated.IsAdmin = userUpdated.IsAdmin
	}
	updated, ok := CheckUpdateUser(user, &userUpdated)
	if !ok {
		WriteJSON(w, http.StatusOK, "Nothing to update.")
		return
	}

	update, err := api.DB.UpdateUser(updated)
	if err != nil {
		api.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, update)
}
func (api *API) Admin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		api.Error(w, http.StatusMethodNotAllowed, "Method not allowed.")
		return
	}
	users, err := api.DB.Users()
	if err != nil {
		api.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, users)
}

// POSTS
func (api *API) CreatePost(w http.ResponseWriter, r *http.Request) {
	var Post models.Post
	if err := json.NewDecoder(r.Body).Decode(&Post); err != nil {

		api.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	Post.UserID = IDUser
	Post.CreatedAt = time.Now().UTC()
	if err := api.DB.CreatePost(&Post); err != nil {
		api.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, &Post)
}

func (api *API) PostByID(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		api.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	if api.UserPosts.Posts[id] == nil {
		api.Error(w, http.StatusNotFound, "No valid user.")
		return
	}
	res, err := api.DB.Post(id)
	if err != nil {
		api.Error(w, http.StatusNotFound, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, res)
}

func (api *API) UpdatePost(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		api.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	if api.UserPosts.Posts[id] == nil {
		api.Error(w, http.StatusNotFound, "No valid user.")
		return
	}
	var postU models.Post
	if err := json.NewDecoder(r.Body).Decode(&postU); err != nil {
		api.Error(w, http.StatusBadRequest, err.Error())
	}
	postU.UpdatedAt = time.Now().UTC()
	post, err := api.DB.Post(id)
	if err != nil {
		api.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := api.DB.UpdatePost(&postU, post); err != nil {
		api.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, postU)
}

func (api *API) DeletePost(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		api.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	if api.UserPosts.Posts[id] == nil {
		api.Error(w, http.StatusNotFound, "No valid user.")
		return
	}
	res, err := api.DB.DeletePost(id)
	if err != nil {
		api.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, res)
}

func (api *API) Posts(w http.ResponseWriter, r *http.Request) {
	res, err := api.DB.Posts(IDUser)
	if err != nil {
		api.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, res)
}
