package main

import (
	"context"
	"encoding/json"
	"fmt"
	_ "github.com/golang-migrate/migrate"
	_ "github.com/lib/pq"
	"groupie-tracker/internal/db"
	"groupie-tracker/internal/models"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	isDataArtistsWritten bool
	pathToCache          = "./internal/cache/cacheArtist.json"
	Mux                  *http.ServeMux
	s                    = "user: -1"
)

type API struct {
	Cache     map[string]*template.Template
	CacheUser *models.User
	UserPosts *models.UsersPosts
	Artists   []models.Artists
	DB        *db.PostgreSQL
}

func Start() {
	logFile, err := os.Create("./internal/logs/result.log")
	if err != nil {
		log.Println("Doesn't open file: ", err, ".")
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetPrefix("Log: ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	postgres, err := OpenDB()
	if err != nil {
		log.Println(err)
	}

	defer postgres.Close()

	context := context.Background()
	ctx, cancel := signal.NotifyContext(context, os.Interrupt)
	defer cancel()

	app := &API{
		Cache: nil,
		DB: &db.PostgreSQL{
			DBSql: postgres,
		},
	}
	if err := InitByMigrateFiles(postgres); err != nil {
		log.Println(err)
	}
	app.checkCacheFile()
	go app.connectingToRedis(ctx)
	go app.checkingInternetConnection(ctx)
	go app.UpdateJSON(ctx)
	go app.Start(ctx, app.Server())
	time.Sleep(5 * time.Second)
	<-ctx.Done()
	time.Sleep(3 * time.Second)
}

func (api *API) connectingToRedis(context context.Context) {
	client := db.RedisDB()
	ticker := time.NewTicker(15 * time.Second)
	i := 1
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			res, err := json.Marshal(api.UserPosts)
			if err != nil {
				log.Println("Can't marshal json for redis: ", err, ".")
			}
			if api.UserPosts != nil {
				s = fmt.Sprintf("user: %v", api.UserPosts.UserP.ID)
			}
			if err := client.Set(context, s, string(res), 5*time.Minute).Err(); err != nil {
				log.Println("Can't insert data to redis:", err, ".")
			}
			log.Println("Data for Redis successfully inserted!!!")
			i++
		case <-context.Done():
			return
		}
	}
}

func (api *API) checkingInternetConnection(context context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	ticks := 0
	for {
		select {
		case t := <-ticker.C:
			if InternetConnection() && !isDataArtistsWritten {
				log.Println(t.Date, " Connecting to the Internet.")
			} else {
				log.Println(t.Date, "Doesn't connecting to the Internet.")
				if ticks%6 == 0 || ticks == 0 {
					isDataArtistsWritten = api.UpdatingCache(api.Artists)
					log.Println(t.Date, " Success update of Artists!!!")
				}
				ticks++
			}
		case <-context.Done():
			isDataArtistsWritten = api.UpdatingCache(api.Artists)
			log.Println("Success update a cache!!!")
			return
		}
	}
}

func InternetConnection() bool {
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(context, "GET", "http://clients3.google.com/generate_204", nil)
	if err != nil {
		return false
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	isDataArtistsWritten = false
	return resp.StatusCode == http.StatusNoContent
}

func (api *API) UpdatingCache(artistsJSON []models.Artists) bool {
	cache, err := json.Marshal(artistsJSON)
	if err != nil {
		log.Println("Marshal json error: ", err, ".")
		return false
	}
	if err := api.CreateCacheFileJSON("cacheArtist.json", cache); err != nil {
		log.Println("Error creating a json: ", err, ".")
		return false
	}
	return true
}

func (api *API) UpdateJSON(context context.Context) {
	go func() {
		for {
			select {
			case <-time.Tick(time.Second*2 + time.Millisecond*500):
				api.UpdateArtistsJSON()
			case <-context.Done():
				return
			}
		}
	}()
}

func (api *API) Start(ctx context.Context, servak *http.Server) {
	go func() {
		<-ctx.Done()

		println("All process stopped.")
		offContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := servak.Shutdown(offContext); err != nil {
			log.Fatal("Can't shutdown: ", err, ".")
		} else {
			println("Have a nice a day!!!")
		}
	}()
	if err := servak.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Error from starting server: ", err, ".")
	}
}

func (api *API) checkCacheFile() {
	if _, err := os.Stat(pathToCache); err == nil {
		data, err := ioutil.ReadFile(pathToCache)
		if err != nil {
			log.Println("Can't a read file: ", err)
		} else {
			if err := json.Unmarshal(data, &api.Artists); err != nil {
				log.Println("Can't unmarshal file: ", err)
			} else {
				log.Println("OK.")
			}
		}
	} else {
		log.Println("Doesn't have any types of this file.")
	}
}
