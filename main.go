package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	port := os.Getenv("PORT")
	token := os.Getenv("GITHUB_TOKEN")
	defaultOrg := os.Getenv("GITHUB_DEFAULT_ORG")
	defaultRepo := os.Getenv("GITHUB_DEFAULT_REPO")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	if token == "" {
		log.Fatal("$GITHUB_TOKEN must be set")
	}
	if defaultOrg == "" {
		log.Fatal("$GITHUB_DEFAULT_ORG must be set")
	}
	if defaultRepo == "" {
		log.Fatal("$GITHUB_DEFAULT_REPO must be set")
	}

	// Establish global github client
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	// Endpoints
	router := gin.New()
	router.Use(gin.Logger())
	router.GET("/stars", func(c *gin.Context) {
		target := c.Query("text")
		segments := strings.Split(target, " ")
		log.Print("input: " + target)
		log.Print("segments:" + strconv.Itoa(len(segments)))

		repo := defaultRepo
		org := defaultOrg

		if len(segments) > 1 {
			org = segments[0]
			repo = segments[1]
		} else if len(segments) > 0 && segments[0] != "" {
			repo = segments[0]
		}

		log.Print("org=" + org + ",repo=" + repo)

		var count string
		repos, _, err := client.Repositories.ListByOrg(org, nil)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to read repositories.")
		}
		for _, r := range repos {
			if *r.Name == repo {
				count = strconv.Itoa(*r.StargazersCount)
				break
			}
		}

		if count != "" {
			slashed := org + "/" + repo
			href := "<https://github.com/" + slashed + "|" + slashed + ">"
			message := href + " has " + count + " :star:"
			c.JSON(http.StatusOK, gin.H{
				"text":          message,
				"response_type": "in_channel",
			})
		} else {
			c.String(http.StatusNotFound, "Repository `"+repo+"` not found.")
		}
	})

	router.Run(":" + port)
}
