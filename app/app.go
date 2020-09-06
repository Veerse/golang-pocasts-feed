package app

import (
	"database/sql"
	"fmt"
	"github.com/Veerse/podcast-feed-api/config"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var (
	LogInfo  *log.Logger
	LogError *log.Logger
)

type App struct {
	Config   config.Config
	Router   *gin.Engine
	DB       sql.DB
	AppCache Cache
}

// Cache is used to store the podcast list and the feeds in memory. The podcast list is the list of all the podcasts
// with their episodes, the feeds are the list of RSS feeds. The key of these two maps is the PodcastID.
// Using a cache tremendously improves performances, going for example from 4.600.000 ns/op to 12.500 ns/op
// for the endpoint /podcasts
type Cache struct {
	Podcasts map[int]Podcast
	Feeds    map[int][]byte
}

func (a *App) Initialize(c config.Config) error {
	a.Config = c

	if err := a.initializeLogger(); err != nil {
		return err
	}
	if err := a.initializeDB(); err != nil {
		LogError.Printf("Initialization: %s\n", err.Error())
		return err
	}
	if err := a.initializeRoutes(); err != nil {
		LogError.Printf("Initialization: %s\n", err.Error())
		return err
	}
	if err := a.initializeCache(); err != nil {
		LogError.Printf("Initialization: %s\n", err.Error())
		return err
	}

	LogInfo.Printf("Initialization successful")
	return nil
}

func (a *App) initializeLogger() error {
	logfile, err := os.OpenFile(a.Config.LogFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.FileMode(0666))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Initializing logger : %s", err.Error())
	}

	LogInfo = log.New(logfile, "Info:\t", log.Ldate|log.Ltime|log.Lshortfile)
	LogError = log.New(logfile, "Error:\t", log.Ldate|log.Ltime|log.Lshortfile)

	return nil
}

func (a *App) initializeDB() error {
	uri := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		a.Config.DB.Host, a.Config.DB.Port, a.Config.DB.User, a.Config.DB.Password, a.Config.DB.Name)

	db, err := sql.Open("postgres", uri)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	a.DB = *db
	return nil
}

func (a *App) initializeRoutes() error {
	a.Router = gin.New()

	a.Router.GET("/podcasts", GetAllPodcasts(&a.AppCache))
	a.Router.GET("/podcasts/:id", GetPodcastById(&a.AppCache))
	a.Router.GET("/podcasts/:id/feed.xml", GetPodcastFeed(&a.AppCache))

	return nil
}

func (a *App) initializeCache() error {
	a.AppCache.Podcasts = make(map[int]Podcast)
	a.AppCache.Feeds = make(map[int][]byte)

	podcasts, err := GetAllPodcastsDao(&a.DB)
	if err != nil {
		return err
	}

	for _, p := range podcasts {
		a.AppCache.Podcasts[p.Id] = p
		a.AppCache.Feeds[p.Id], _ = p.ToFeed()
	}

	return nil
}

func (a *App) Run() {
	//gin.SetMode(gin.ReleaseMode)
	LogInfo.Printf("Starting server")
	a.Router.Run()
}
