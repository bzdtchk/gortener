package http

import (
	"Gortener/database"
	sqlDb "database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var db *sqlDb.DB

type shortenRequest struct {
	DestinationUrl string `json:"destination_url" binding:"required"`
	Slug           string `json:"slug" binding:"required"`
}

type shortenResponse struct {
	Slug string `json:"slug"`
}

func Run() {
	ginEngine := gin.Default()
	ginEngine.LoadHTMLFiles("./www/index.html")
	ginEngine.StaticFile("/shorten.js", "./www/shorten.js")
	setupRouter(ginEngine)

	ginEngine.Run()
}

func setupRouter(ginEngine *gin.Engine) {
	ginEngine.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	ginEngine.POST("/api/shorten-link", func(c *gin.Context) {
		var shortenRequest shortenRequest

		if err := c.BindJSON(&shortenRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid JSON received",
			})
			return
		}

		isAlphanumericalOnly, _ := regexp.MatchString(`^[a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)*$`, shortenRequest.Slug)
		if !isAlphanumericalOnly {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Slug must be an alphanumerical string",
			})
			return
		}
		shortenRequest.Slug = strings.ToLower(shortenRequest.Slug)

		link, err := database.GetBySlug(shortenRequest.Slug)
		if link != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "There's already link with the same slug in our database",
			})
			return
		}

		req, err := http.NewRequest("GET", shortenRequest.DestinationUrl, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Destination URL for shortening is not reachable",
			})
			return
		}
		defer resp.Body.Close()

		link, err = database.InsertLink(&database.ShortenedLink{Slug: shortenRequest.Slug, DestinationUrl: shortenRequest.DestinationUrl})
		if err != nil {
			return
		}

		c.IndentedJSON(http.StatusCreated, shortenResponse{Slug: link.Slug})
	})

	ginEngine.GET("/:slug", func(c *gin.Context) {
		slugRequestValue := c.Param("slug")

		link, err := database.GetBySlug(slugRequestValue)
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.Redirect(http.StatusMovedPermanently, link.DestinationUrl)
	})
}
