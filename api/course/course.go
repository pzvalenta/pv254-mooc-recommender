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
	ID               string
	Categories       []string
	Description      string
	Details          CourseDetails
	Interested_Count string
	Link             string
	Name             string
	Overview         string
	Provider         string
	Rating           int32
	Review_Count     string
	Schools          []string
	Subject          string
	Syllabus         string
	Teachers         []string
}

// CourseDetails ...
type CourseDetails struct {
	Certificate        string
	Cost               int32
	Currency           string
	Duration           int32
	Duration_Time_Unit string
	Effort             int32
	Effort_Time_Unit   string
	Language           string
	Provider           string
	Session            string
	StartDate          []string
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
