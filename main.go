package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"os"
	"time"
)

var (
	httpClient *http.Client
	baseURL string
)

func main(){
	httpClient = &http.Client{}
	baseURL = "https://api.ldjam.com"

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.LoadHTMLGlob("/app/public/templates/*.html")
	r.GET("/", showPage)
	r.POST("/", getRank)
	r.Static("/assets", "./public")
	r.Run(":"+getEnv("PORT", "8080"))
}

func showPage(c *gin.Context){
	c.HTML(http.StatusOK, "index.html", map[string]interface{}{
		"Test": time.Date(2017, 07, 01, 0, 0, 0, 0, time.UTC),
	})
}

func getRank(c *gin.Context){
	var game *LDGame
	var event *LDEvent
	var message string

	client := NewClient(baseURL, httpClient)
	userInput := c.PostForm("url")
	userURL, err := url.ParseRequestURI(userInput)

	if err != nil{
		message = "Invalid URL"
	} else {
		game, err = client.GetGameFromURL(userURL)
		event, err = client.GetEventStatsFromGame(game)
	}

	c.HTML(http.StatusOK, "index.html", map[string]interface{}{
		"input": userInput,
		"game": game,
		"event": event,
		"message": message,
	})

}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}