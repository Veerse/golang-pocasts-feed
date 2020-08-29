package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Veerse/podcast-feed-api/config"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func setupLoggers() {
	logfile, err := os.OpenFile("../logs_test.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, os.FileMode(0666))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Initializing logger : %s", err.Error())
	}

	LogInfo = log.New(logfile, "Info:\t", log.Ldate|log.Ltime|log.Lshortfile)
	LogError = log.New(logfile, "Error:\t", log.Ldate|log.Ltime|log.Lshortfile)
}

func setupConfig() config.Config {
	c := config.Config{}

	// As we are on the /app folder, the configapp.json is located on level above
	if f, err := os.Open("../configapp.json"); err != nil {
		panic(err)
	} else {
		defer f.Close()
		json.NewDecoder(f).Decode(&c)
		return c
	}
}

func setupDB(c config.Config) *sql.DB {
	uri := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.DB.Host, c.DB.Port, c.DB.User, c.DB.Password, c.DB.Name)

	db, err := sql.Open("postgres", uri)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	return db
}

func setupCache(db *sql.DB) Cache {
	cache := Cache{}

	cache.Podcasts = make(map[int]Podcast)
	cache.Feeds = make(map[int]string)

	podcasts, err := GetAllPodcastsDao(db)
	if err != nil {
		panic(err)
	}

	for _, p := range podcasts {
		cache.Podcasts[p.Id] = p
	}

	return cache
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	c := setupConfig()

	db := setupDB(c)

	cache := setupCache(db)

	router.GET("/podcasts", GetAllPodcasts(&cache))
	router.GET("/podcasts/:id", GetPodcastById(&cache))

	return router
}

func BenchmarkGetAllPodcasts(b *testing.B) {
	router := setupRouter()

	gin.DefaultWriter = ioutil.Discard

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/podcasts", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}

func BenchmarkGetPodcastById(b *testing.B) {
	setupLoggers()

	router := setupRouter()

	gin.DefaultWriter = ioutil.Discard

	w := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "/podcasts/1", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}

func TestGetAllPodcasts(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/podcasts", nil)

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Unexpected HTTP code %d", w.Code)
	}
}

func TestGetPodcastById(t *testing.T) {
	setupLoggers()
	router := setupRouter()

	w := httptest.NewRecorder()

	for i := 0; i <= 5; i++ {
		uri := fmt.Sprintf("/podcasts/%d", i)
		req := httptest.NewRequest("GET", uri, nil)
		router.ServeHTTP(w, req)

		if w.Code != 200 && w.Code != 204 {
			t.Errorf("Unexpected HTTP code %d", w.Code)
		}
	}
}
