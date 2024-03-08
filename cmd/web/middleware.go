package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"os"
	"time"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))
var IDUser = -1
var IsAdmin = false

func (api *API) GenerateJWT(email string, is_admin bool, name string, id int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["id"] = id
	claims["email"] = email
	claims["isAdmin"] = is_admin
	claims["name"] = name
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("error in GenerateJWT: %w", err)
	}
	return tokenString, nil
}

func (api *API) IsAuthorized(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		token, err := AuthChecker(authHeader)
		if err != nil {
			api.Error(w, http.StatusUnauthorized, err.Error())
			return
		}
		IDUser = int(int64(token.Claims.(jwt.MapClaims)["id"].(float64)))
		next(w, r)
	})
}

func (api *API) IsAuthorizedJWT(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		token, err := AuthChecker(authHeader)
		if err != nil {
			api.Error(w, http.StatusUnauthorized, err.Error())
			return
		}
		id, err := getID(r)
		claims := token.Claims.(jwt.MapClaims)
		if id != int(int64(claims["id"].(float64))) && !claims["isAdmin"].(bool) {
			http.Error(w, "No Valid User", http.StatusBadRequest)
			return
		}
		if claims["isAdmin"].(bool) {
			IsAdmin = true
		}
		next(w, r)
	})
}

//	if id != int(int64(claims["id"].(float64)))   {
//		http.Error(w, "No Valid User", http.StatusBadRequest)
//		return
//	}
func (api *API) IsAdmin(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		token, err := AuthChecker(authHeader)
		if err != nil {
			api.Error(w, http.StatusUnauthorized, err.Error())
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		if !claims["isAdmin"].(bool) {
			log.Println("Access rejected because you are not admin")
			http.Error(w, "Access rejected because you are not admin", http.StatusUnauthorized)
			return
		}
		IsAdmin = true
		next(w, r)
	})
}

func enableCorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
