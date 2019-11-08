package main

import (
	"context"
	"net/http"
	"time"

	"log"

	"go-docker-api/api/db"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	router := gin.Default()
	public := router.Group("/api")
	{
		public.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "Welcome to sample dockerized golang api")
		})
	}
	time.Sleep(20 * time.Second) //sleep until db image is running

	c := db.GetClient()
	err := c.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Couldn't connect to the database", err)
	} else {
		log.Println("Connected!")
	}

	router.Run()
}
