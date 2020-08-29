package app

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllPodcasts(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if podcasts, err := GetAllPodcastsDao(db); err != nil {
			LogError.Printf("GetAllPodcasts : %s", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		} else {
			c.JSON(http.StatusOK, podcasts)
		}
	}
}

func GetPodcastById(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		podcastId := c.Param("id")

		if podcast, err := GetPodcastByIdDao(podcastId, db); err != nil {
			if err == sql.ErrNoRows {
				c.Status(http.StatusNoContent)
			} else {
				LogError.Printf("GetPodcastById : %s", err.Error())
				c.Status(http.StatusInternalServerError)
			}
		} else {
			c.JSON(http.StatusOK, podcast)
		}
	}
}

func GetPodcastFeed(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Data(200, "text/plain", []byte("Work in progress ;)"))
	}
}
