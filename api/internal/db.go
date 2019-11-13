package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewDatabase(host, port string) (*mongo.Database, error) {
	if envHost, ok := os.LookupEnv("DB_HOST"); ok {
		host = envHost
	}

	clientOptions := options.Client().ApplyURI("mongodb://" + host + ":" + port)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %v", err)
	}
	err = client.Connect(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error connecting to client: %v", err)
	}

	return client.Database("mydb"), nil
}

func (s *State) RandomCourse(c *gin.Context) {
	coursesCollection := s.DB.Collection("courses")
	var result Course

	filter := bson.D{{}} // empty filter
	err := coursesCollection.FindOne(c, filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, result)
}

func (s *State) TaxonomyCourses(c *gin.Context) {
	myCourseIds, err := s.getMyCoursesIds(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	myCourses, err := s.getMyCourses(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
	}

	myCoursesSubjects := ExtractSubjects(c, myCourses)
	coursesCollection := s.DB.Collection("courses")
	filter := bson.D{
		{Key: "_id", Value: bson.D{{Key: "$nin", Value: myCourseIds}}},
		{Key: "subject", Value: bson.D{{Key: "$in", Value: myCoursesSubjects}}},
		//{"categories", bson.D{{"$nin", bson.A{myCourseCategories}}}}
	}
	data, err := coursesCollection.Find(c, filter, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
	}

	var coursesFromOtherSubtree []Course
	for data.Next(c) {
		l := Course{}
		err = data.Decode(&l)
		if err != nil {
			c.JSON(http.StatusInternalServerError, nil)
		}
		coursesFromOtherSubtree = append(coursesFromOtherSubtree, l)
	}

	similar := FindSimilar(c, myCourses, coursesFromOtherSubtree)[:10]
	c.JSON(http.StatusOK, similar)
}

func (s *State) getMyCoursesIds(c *gin.Context) ([]string, error) {
	users := s.DB.Collection("users")

	id, err := primitive.ObjectIDFromHex(s.customerID)
	if err != nil {
		return nil, fmt.Errorf("error creating id from hex: %v", err)
	}

	filter := bson.M{"_id": id}
	data, err := users.Find(c, filter)
	if err != nil {
		return nil, fmt.Errorf("unable to find user's course IDs: %v", err)
	}

	l := User{}
	if data.Next(c) {
		err = data.Decode(&l)
		if err != nil {
			return nil, fmt.Errorf("unable to decode user's enrolledIn: %v", err)
		}
	}

	return l.EnrolledIn, nil
}

func (s *State) getMyCourses(c *gin.Context) ([]Course, error) {
	myCourseIDs, err := s.getMyCoursesIds(c)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	ids := bson.A{}
	for _, id := range myCourseIDs {
		ids = append(ids, id)
	}

	var res []Course
	filter := bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: ids}}}}
	courseCollection := s.DB.Collection("courses")
	data, err := courseCollection.Find(c, filter, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to find/decode courses %v", err)
	}

	for data.Next(c) {
		l := Course{}
		err = data.Decode(&l)
		if err != nil {
			return nil, fmt.Errorf("unable to decode user's enrolledIn: %v", err)
		}
		res = append(res, l)
	}

	return res, nil
}