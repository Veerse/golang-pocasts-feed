package app

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllPodcasts(db *sql.DB) gin.HandlerFunc {
	return func (c *gin.Context) {
		if podcasts, err := GetAllPodcastsDao(db); err != nil {
			LogError.Printf("GetAllPodcasts : %s", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		} else {
			fmt.Println("La len des podcasts est de", len(podcasts))
			c.JSON(http.StatusOK, podcasts)
		}
	}
}

func GetPodcastById(db *sql.DB) gin.HandlerFunc {
	return func (c *gin.Context) {
		podcastId := c.Param("id")

		if podcast, err := GetPodcastByIdDao(podcastId, db); err != nil {
			LogError.Printf("GetPodcastById : %s", err.Error())
			if err == sql.ErrNoRows {
				c.Status(http.StatusNoContent)
			} else {
				c.Status(http.StatusInternalServerError)
			}
		} else {
			c.JSON(http.StatusOK, podcast)
		}
	}
}