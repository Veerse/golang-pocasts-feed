package app

import (
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
		id := c.Param("id")

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
		id := c.Param("id")

		podcastId, _ := strconv.Atoi(id)

		f, exists := cache.Feeds[podcastId]

		if !exists {
			c.Status(http.StatusNotFound)
		} else {
			c.Data(http.StatusOK, "text/xml", f)
		}
	}
}
