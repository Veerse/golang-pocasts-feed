package app

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Veerse/podcast-feed-api/config"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	LogInfo  *log.Logger
	LogError *log.Logger
)

type App struct {
	Config         config.Config
	Router         *gin.Engine
	DB             sql.DB
	AppCache       Cache
	AuthMiddleware *jwt.GinJWTMiddleware
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
		return errors.New(fmt.Sprintf("initialize logger : %s", err.Error()))
	}
	if err := a.initializeRouter(); err != nil {
		LogError.Printf("initialize router: %s", err.Error())
		return errors.New(fmt.Sprintf("initialize router : %s", err.Error()))
	}
	if err := a.initializeJWT(); err != nil {
		LogError.Printf("initialize JWT: %s", err.Error())
		return errors.New(fmt.Sprintf("initialize JWT : %s", err.Error()))
	}
	if err := a.initializeDB(); err != nil {
		LogError.Printf("initialize db: %s", err.Error())
		return errors.New(fmt.Sprintf("initialize db : %s", err.Error()))
	}
	if err := a.initializeRoutes(); err != nil {
		LogError.Printf("initialize routes: %s", err.Error())
		return errors.New(fmt.Sprintf("initialize routes : %s", err.Error()))
	}
	if err := a.initializeCache(); err != nil {
		LogError.Printf("initialize cache: %s", err.Error())
		return errors.New(fmt.Sprintf("initialize cache : %s", err.Error()))
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

func (a *App) initializeRouter() error {
	a.Router = gin.New()
	a.Router.Use(cors.Default())

	return nil
}

func (a *App) initializeJWT() error {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "Muslimy",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: "id",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					"id":        v.Id,
					"privilege": v.Privilege,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				Id:        int(claims["id"].(float64)),
				Privilege: int(claims["privilege"].(float64)),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			type login struct {
				Email    string `form:"email" json:"email" binding:"required"`
				Password string `form:"password" json:"password" binding:"required"`
			}

			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			if u, err := GetUserByEmailAndPassword(loginVals.Email, loginVals.Password, &a.DB); err == nil {
				return &u, nil
			} else {
				if !errors.Is(err, sql.ErrNoRows) {
					LogError.Printf("authentificator authentification : %s", err.Error())
					return "", err
				}
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if u, ok := data.(*User); ok {
				if u.Privilege == unverified {
					return false
				}

				if u.Privilege == admin {
					return true
				}

				if id := c.Param("podcastId"); id != "" {
					podcastId, _ := strconv.Atoi(id)
					if f, exists := a.AppCache.Podcasts[podcastId]; exists {
						if f.UserId == u.Id && (u.Privilege == poster) {
							return true
						}
					}
				}
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization", //, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		return err
	}

	if errInit := authMiddleware.MiddlewareInit(); errInit != nil {
		return err
	}

	a.AuthMiddleware = authMiddleware
	return nil
}

func (a *App) initializeRoutes() error {
	a.Router.POST("/login", a.AuthMiddleware.LoginHandler)

	a.Router.GET("/podcasts", GetAllPodcasts(&a.AppCache))
	a.Router.GET("/podcasts/:podcastId", GetPodcastById(&a.AppCache))
	a.Router.GET("/podcasts/:podcastId/feed.xml", GetPodcastFeed(&a.AppCache))

	auth := a.Router.Group("/")

	auth.Use(a.AuthMiddleware.MiddlewareFunc())
	{
		auth.GET("/hello", HelloHandler)
		a.Router.POST("/podcasts", CreatePodcast())
		a.Router.POST("/podcasts/:podcastId/episodes", CreateEpisode(&a.AppCache))
	}

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
