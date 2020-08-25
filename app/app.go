package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Veerse/podcast-feed-api/config"
	"github.com/eduncan911/podcast"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	LogInfo		*log.Logger
	LogError	*log.Logger
)
var Podcasts []podcast.Podcast
var Podcasts2 []feeds.Feed


type App struct {
	Config		config.Config
	Router		*gin.Engine
	DB			*sql.DB
	Feeds		map[int]string
}

func (a *App) Initialize (c config.Config) error {
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
	if err := a.initializeFeeds(); err != nil {
		LogError.Printf("Initialization: %s\n", err.Error())
		return err
	}
	if err := a.initializePodcasts(); err != nil {
		LogError.Printf("Initialization: %s\n", err.Error())
		return err
	}
	if err := a.initializePodcasts2(); err != nil {
		LogError.Printf("Initialization: %s\n", err.Error())
		return err
	}

	LogInfo.Printf("Initialization successful")
	return nil
}

func (a *App) initializeLogger () error {
	logfile, err := os.OpenFile(a.Config.LogFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.FileMode(0666))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Initializing logger : %s", err.Error())
	}

	LogInfo = log.New(logfile, "Info:\t", log.Ldate|log.Ltime|log.Lshortfile)
	LogError = log.New(logfile, "Error:\t", log.Ldate|log.Ltime|log.Lshortfile)

	return nil
}

func (a *App) initializeDB () error {
	uri := "host=localhost port=5432 user=postgres "+
		"password=root dbname=postgres sslmode=disable sslmode=disable"

	db, err := sql.Open("postgres", uri)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}
	fmt.Printf("DB initialization complete \n")
	return nil
}

func (a *App) initializeRoutes () error {
	a.Router = gin.New()

	a.Router.GET("/", myFunc)
	a.Router.GET("/2", myFunc2)
	return nil
}

func (a *App) initializeFeeds () error {
	return nil
}

// EDUNCAN
func (a *App) initializePodcasts () error {
	p := podcast.New(
		"eduncan911 Podcasts",
		"http://eduncan911.com/",
		"An example Podcast",
		nil, nil,
	)
	p.Language = "fr-FR"
	p.AddAuthor("Jane Doe", "me@janedoe.com")
	p.AddImage("http://janedoe.com/i.jpg")
	p.AddSummary(`A very cool podcast nilwith a long summary using Bytes()!
See more at our website: <a href="http://example.com">example.com</a>
`)
	p.IExplicit = "no"
	p.Category = "Religion &amp; Spirituality"

	for i := int64(5); i < 7; i++ {
		n := strconv.FormatInt(i, 10)

		item := podcast.Item{
			Title:       "Episode " + n,
			Link:        "http://example.com/" + n + ".mp3",
			Description: "Description for Episode " + n,
			PubDate:     nil,
		}
		if _, err := p.AddItem(item); err != nil {
			fmt.Println(item.Title, ": error", err.Error())
			break
		}
	}

	fmt.Printf("Len : %d\n", len(Podcasts))
	Podcasts = append(Podcasts, p)
	fmt.Printf("Len : %d\n", len(Podcasts))
	return nil
}

// GORILLA
func (a *App) initializePodcasts2 () error {
	now := time.Now()
	feed := feeds.Feed{
		Title:       "jmoiron.net blog",
		Link:        &feeds.Link{Href: "http://jmoiron.net/blog"},
		Description: "discussion about tech, footie, photos",
		Author:      &feeds.Author{Name: "Jason Moiron", Email: "jmoiron@jmoiron.net"},
		Created:     now,
	}

	feed.Items = []*feeds.Item{
		&feeds.Item{
			Title:       "Limiting Concurrency in Go",
			Link:        &feeds.Link{Href: "http://jmoiron.net/blog/limiting-concurrency-in-go/"},
			Description: "A discussion on controlled parallelism in golang",
			Author:      &feeds.Author{Name: "Jason Moiron", Email: "jmoiron@jmoiron.net"},
			Created:     now,
		},
		&feeds.Item{
			Title:       "Logic-less Template Redux",
			Link:        &feeds.Link{Href: "http://jmoiron.net/blog/logicless-template-redux/"},
			Description: "More thoughts on logicless templates",
			Created:     now,
		},
		&feeds.Item{
			Title:       "Idiomatic Code Reuse in Go",
			Link:        &feeds.Link{Href: "http://jmoiron.net/blog/idiomatic-code-reuse-in-go/"},
			Description: "How to use interfaces <em>effectively</em>",
			Created:     now,
		},
	}

	Podcasts2 = append(Podcasts2, feed)
	fmt.Printf("la len mon gars %d", len(Podcasts2))
	return nil
}

func (a *App) Run () {
	LogInfo.Printf("Starting server")
	a.Router.Run()
}

func myFunc (c *gin.Context) {
	b := Podcasts[0]

	c.Header("content-type", "application/json")
	json.NewEncoder(c.Writer).Encode(b)
	//c.JSON(200, Podcasts)

	//c.Data(200, "application/xml", p.Bytes())
}

func myFunc2 (c *gin.Context) {
	b, _ := Podcasts2[0].ToJSON()
	fmt.Printf("La len %d", len(Podcasts2))
	c.Data(200, "application/json", []byte(b))
}