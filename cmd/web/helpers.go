package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"groupie-tracker/internal/models"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	api         = "https://groupietrackers.herokuapp.com/api/artists"
	apiRelation = "https://groupietrackers.herokuapp.com/api/relation"
)

func Artists() ([]models.Artists, error) {
	var artists []models.Artists
	var relation models.Relations

	artistsApi, err := http.Get(api)
	if err != nil {
		return nil, err
	}
	defer artistsApi.Body.Close()

	err = json.NewDecoder(artistsApi.Body).Decode(&artists)
	if err != nil {
		return nil, err
	}

	relationApi, err := http.Get(apiRelation)
	if err != nil {
		return nil, err
	}
	defer relationApi.Body.Close()

	err = json.NewDecoder(relationApi.Body).Decode(&relation)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(artists); i++ {
		artists[i].DateLocations = relation.Index[i].DateLocations
	}
	return artists, nil
}

var (
	locker  sync.RWMutex
	artists []models.Artists
)

func (api *API) CreateCacheFileJSON(fileName string, dataJSON []byte) error { // api ???
	file, err := os.Create("./internal/cache/" + fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(dataJSON)
	if err != nil {
		return err
	}
	return nil
}

func (api *API) UpdateArtistsJSON() error {
	updatedArtists, err := Artists()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	locker.Lock()
	defer locker.Unlock()
	api.Artists = updatedArtists
	log.Println("Cache updated.")
	return nil
}

func (api *API) Error(w http.ResponseWriter, statusCode int, errorMessage string) {
	log.Println(errorMessage)
	w.WriteHeader(statusCode)
	err := &models.Error{ErrorMessage: errorMessage, StatusCode: statusCode}
	WriteJSON(w, statusCode, err)
}
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

// GORILLA
func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("ID %s: is invalid", idStr)
	}
	return id, nil
}

// DB
func OpenDB() (*sql.DB, error) {
	connStr := "host=localhost port=5432 user=postgres password=0000 dbname=golang sslmode=disable"
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	log.Println("Connected Successfully: Postgres")
	if err := dbConn.Ping(); err != nil {
		return nil, err
	}
	return dbConn, nil
}

// USE A MIGRATION
func InitByMigrateFiles(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err

	}
	m, err := migrate.NewWithDatabaseInstance("file:///Users/zhalgasbolatov/Downloads/groupie-tracker-6/internal/migrations", "postgres", driver)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

// FOR MIDDLEWARE
func AuthChecker(authHeader string) (*jwt.Token, error) {
	if authHeader == "" {
		log.Println("Authentication header is required")
		return nil, fmt.Errorf("Authentication header is require")
	}
	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 {
		log.Println("Invalid Authorization token format")
		return nil, fmt.Errorf("Invalid Authorization token format")
	}
	token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if !token.Valid {
		log.Println("Invalid token")
		return nil, fmt.Errorf("Invalid token")
	}
	return token, nil
}

func MapOfPosts(postsArr []*models.Post) map[int]*models.Post {
	posts := map[int]*models.Post{}
	for _, v := range postsArr {
		posts[v.ID] = v
	}
	return posts
}
func CheckUpdateUser(user, userUpdated *models.User) (*models.User, bool) {
	check := false
	if user.Name != userUpdated.Name {
		user.Name = userUpdated.Name
		check = true
	}
	if user.Email != userUpdated.Email {
		user.Email = userUpdated.Email
		check = true
	}
	if userUpdated.HashPassword != "" {
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(userUpdated.HashPassword), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(err)
		}
		user.HashPassword = string(hashPassword)
	}
	if user.IsAdmin != userUpdated.IsAdmin {
		user.IsAdmin = userUpdated.IsAdmin
		check = true
	}

	return user, check
}
