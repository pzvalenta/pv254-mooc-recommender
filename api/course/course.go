package course

import (
	"context"
	"go-docker-api/api/db"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// Course ...
type Course struct {
	ID              string   `bson:"_id"`
	Categories      []string `bson:"categories"`
	Description     string   `bson:"description"`
	Details         Details  `bson:"details"`
	InterestedCount string   `bson:"interested_count"`
	Link            string   `bson:"link"`
	Name            string   `bson:"name"`
	Overview        string   `bson:"overview"`
	Provider        string   `bson:"provider"`
	Rating          int32    `bson:"rating"`
	ReviewCount     string   `bson:"review_count"`
	Schools         []string `bson:"schools"`
	Subject         string   `bson:"subject"`
	Syllabus        string   `bson:"syllabus"`
	Teachers        []string `bson:"teachers"`
}

// Details ...
type Details struct {
	Certificate      string   `bson:"certificate"`
	Cost             int32    `bson:"cost"`
	Currency         string   `bson:"currency"`
	Duration         int32    `bson:"duration"`
	DurationTimeUnit string   `bson:"duration_time_unit"`
	Effort           int32    `bson:"effort"`
	EffortTimeUnit   string   `bson:"effort_time_unit"`
	Language         string   `bson:"language"`
	Provider         string   `bson:"provider"`
	Session          string   `bson:"session"`
	StartDate        []string `bson:"start date"`
}

// RandomCourse ...
func RandomCourse(c *gin.Context) {
	coursesCollection := db.GetDB().Collection("courses")
	var result Course

	filter := bson.D{{}} // empty filter
	err := coursesCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Found a single document") //: %s\n", result.ID
	c.JSON(http.StatusOK, result)
}
