package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"time"

	"shorty/shortener"

	"github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
)

//go:embed static/*
var staticFS embed.FS

type ShortenRequest struct {
	LongURL string `json:"longURL"`
}

func setupRouter(shorty *shortener.RedisURLShortener, domain string) *gin.Engine {
	r := gin.Default()

	store := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
		Rate:  time.Second,
		Limit: 5,
	})
	mw := ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc:      keyFunc,
	})

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/shorty")
	})

	r.GET("/r/:code", func(c *gin.Context) {
		shortCode := c.Param("code")
		longURL, err := shorty.Redirect(shortCode)
		if err != nil {
			fmt.Println("error: short URL not found")
			c.String(http.StatusNotFound, "Short URL not found")
			return
		}
		c.Redirect(http.StatusFound, longURL)
	})

	staticFS, _ := fs.Sub(staticFS, "static")
	r.StaticFS("/shorty", http.FS(staticFS))

	r.POST("/shorty", mw, func(c *gin.Context) {
		var request ShortenRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			fmt.Println("error: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if !isValidURL(request.LongURL) {
			fmt.Println("error: invalid URL format")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL format"})
			return
		}

		fmt.Println("longURL: ", request.LongURL)
		slug, err := shorty.ShortenURL(request.LongURL)
		if err != nil {
			fmt.Println("error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error: ": err})
			return
		}

		address := fmt.Sprintf("http://%s/r/%s", domain, slug)
		c.JSON(http.StatusOK, gin.H{"message": address})
		return
	})

	return r
}

func errorHandler(c *gin.Context, info ratelimit.Info) {
	c.String(429, "Too many requests. Try again in "+time.Until(info.ResetTime).String())
}

func keyFunc(c *gin.Context) string {
	return c.ClientIP()
}

func isValidURL(u string) bool {
	// URL validation logic here
	return true
}
