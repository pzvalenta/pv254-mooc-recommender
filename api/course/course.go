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
	ID              string   `json:"id" bson:"_id"`
	Categories      []string `json:"categories" bson:"categories"`
	Description     string   `json:"description" bson:"description"`
	Details         Details  `json:"details" bson:"details"`
	InterestedCount string   `json:"interested_count" bson:"interested_count"`
	Link            string   `json:"link" bson:"link"`
	Name            string   `json:"name" bson:"name"`
	Overview        string   `json:"overview" bson:"overview"`
	Provider        string   `json:"provider" bson:"provider"`
	Rating          *float64 `json:"rating" bson:"rating"`
	ReviewCount     string   `json:"review_count" bson:"review_count"`
	Schools         []string `json:"schools" bson:"schools"`
	Subject         string   `json:"subject" bson:"subject"`
	Syllabus        string   `json:"syllabus" bson:"syllabus"`
	Teachers        []string `json:"teachers" bson:"teachers"`
}

// Details ...
type Details struct {
	Certificate      string   `json:"certificate" bson:"certificate"`
	Cost             int32    `json:"cost" bson:"cost"`
	Currency         string   `json:"currency" bson:"currency"`
	Duration         *float64 `json:"duration" bson:"duration"`
	DurationTimeUnit string   `json:"duration_time_unit" bson:"duration_time_unit"`
	Effort           *float64 `json:"effort" bson:"effort"`
	EffortTimeUnit   string   `json:"effort_time_unit" bson:"effort_time_unit"`
	Language         string   `json:"language" bson:"language"`
	Provider         string   `json:"provider" bson:"provider"`
	Session          string   `json:"session" bson:"session"`
	StartDate        []string `json:"start date" bson:"start date"`
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
