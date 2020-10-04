package app

import (
	"database/sql"
	"fmt"
	"github.com/Veerse/podcast-feed-api/config"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"os"
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
		return err
	}
	if err := a.initializeRouter(); err != nil {
		LogError.Printf("Initialization: %s", err.Error())
		return err
	}
	if err := a.initializeJWT(); err != nil {
		LogError.Printf("Initialization: %s", err.Error())
		return err
	}
	if err := a.initializeDB(); err != nil {
		LogError.Printf("Initialization: %s", err.Error())
		return err
	}
	if err := a.initializeRoutes(); err != nil {
		LogError.Printf("Initialization: %s", err.Error())
		return err
	}
	if err := a.initializeCache(); err != nil {
		LogError.Printf("Initialization: %s", err.Error())
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

func (a *App) initializeRouter() error {
	a.Router = gin.New()
	a.Router.Use(cors.Default())

	return nil
}

// User demo
type User struct {
	ID        string
	UserName  string
	FirstName string
	LastName  string
	Age       int
}

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

var identityKey = "id"

func (a *App) initializeJWT() error {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			fmt.Printf("datas %+v\n", data)
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.ID,
					"role":      v.Age,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				ID: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password

			if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
				return &User{
					ID:        userID,
					UserName:  "Test",
					LastName:  "Bo-Yi",
					FirstName: "Wu",
					Age:       145,
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			fmt.Printf("%+v", data)
			if v, ok := data.(*User); ok && v.ID == "admin" {
				claims := jwt.ExtractClaims(c)
				//user, _ := c.Get(identityKey)
				fmt.Printf("%+v", claims)
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatalf("Error %s", err.Error())
	}

	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatalf("Error %s", err.Error())
	}

	a.Router.POST("/login", authMiddleware.LoginHandler)

	auth := a.Router.Group("/auth")
	// Refresh time can be longer than token timeout
	//auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/hello", HelloHandler)
	}

	//a.Router.GET("/hello", HelloHandler).Use(authMiddleware.MiddlewareFunc())
	a.AuthMiddleware = authMiddleware

	return nil
}

func (a *App) initializeRoutes() error {
	a.Router.GET("/podcasts", GetAllPodcasts(&a.AppCache))
	a.Router.GET("/podcasts/:id", GetPodcastById(&a.AppCache))
	a.Router.GET("/podcasts/:id/feed.xml", GetPodcastFeed(&a.AppCache))

	a.Router.POST("/podcasts", CreatePodcast())
	a.Router.POST("/podcasts/:id/episodes", CreateEpisode(&a.AppCache))

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
