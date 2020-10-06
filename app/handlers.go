package app

import (
	"database/sql"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetAllPodcasts(cache *Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		var podcasts []Podcast

		for _, v := range cache.Podcasts {
			podcasts = append(podcasts, v)
		}
		c.JSON(http.StatusOK, podcasts)
	}
}

func GetPodcastById(cache *Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("podcastId")

		podcastId, _ := strconv.Atoi(id)

		p, exists := cache.Podcasts[podcastId]

		if !exists {
			c.Status(http.StatusNotFound)
		} else {
			c.JSON(http.StatusOK, p)
		}
	}
}

func GetPodcastFeed(cache *Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("podcastId")

		podcastId, _ := strconv.Atoi(id)

		f, exists := cache.Feeds[podcastId]

		if !exists {
			c.Status(http.StatusNotFound)
		} else {
			c.Data(http.StatusOK, "text/xml", f)
		}
	}
}

func CreatePodcast() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input Podcast

		if err := c.ShouldBindJSON(input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	}
}

func CreateEpisode(db *sql.DB, cache *Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("podcastId")
		pid, _ := strconv.Atoi(id)

		if _, exists := cache.Podcasts[pid]; exists != true {
			c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "podcast doesn't exist"})
			return
		}

		var input Episode

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "fields missing or incorrect"})
			return
		}

		if input.PodcastId != 0 && input.PodcastId != pid {
			claims := jwt.ExtractClaims(c)
			LogAlert.Printf("user %d tried to add episode to an unowned podcast %d", int(claims["id"].(float64)), input.PodcastId)
			c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "you are not allowed to edit this resource"})
			return
		}

		if e, err := CreateEpisodeDAO(input, db); err != nil {
			LogError.Printf("creating episode DAO : %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "error": "internal error"})
			return
		} else {
			c.JSON(http.StatusCreated, e)
		}
	}
}

func HelloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get("id")
	fmt.Printf("usr : %+v\n", user)
	c.JSON(200, gin.H{
		"userID":    claims["id"],
		"userName":  user.(*User).Id,
		"text":      "Hello World.",
		"privilege": claims["privilege"],
	})
}
