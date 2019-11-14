package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"

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
		c.JSON(http.StatusInternalServerError, "no content")
		return
	}

	myCourses, err := s.getMyCourses(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "no content")
		return
	}

	recommended := make(map[string][]Course)
	coursesCollection := s.DB.Collection("courses")
	for i := range myCourses {
		var courseCategories bson.A
		for _, cat := range myCourses[i].Categories {
			courseCategories = append(courseCategories, cat)
		}
		filter := bson.D{
			{Key: "_id", Value: bson.D{{Key: "$nin", Value: myCourseIds}}},
			{Key: "subject", Value: myCourses[i].Subject},
			{Key: "categories", Value: bson.D{{Key: "$nin", Value: courseCategories}}},
		}

		coursesFromOtherSubtree, err := s.FindCoursesAccordingFilter(c, filter, coursesCollection)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "no content")
			return
		}

		similar := myCourses[i].FindSimilar(coursesFromOtherSubtree)
		sort.Sort(BySimilarity{courses: similar, course: &myCourses[i]})
		recommended[myCourses[i].ID] = similar[:10]
	}

	c.JSON(http.StatusOK, convertResult(recommended))
}

func (s *State) OverfittingCourses(c *gin.Context) {
	myCourseIds, err := s.getMyCoursesIds(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "no content")
		return
	}

	myCourses, err := s.getMyCourses(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "no content")
		return
	}

	recommended := make(map[string][]Course)
	coursesCollection := s.DB.Collection("courses")
	for i := range myCourses {
		filter := bson.D{
			{Key: "_id", Value: bson.D{{Key: "$nin", Value: myCourseIds}}},
			{Key: "categories", Value: bson.D{{Key: "$nin", Value: myCourses[i].Categories}}},
		}
		coursesWithoutMine, err := s.FindCoursesAccordingFilter(c, filter, coursesCollection)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "no content")
			return
		}

		sort.Sort(BySimilarity{course: &myCourses[i], courses: coursesWithoutMine})
		recommended[myCourses[i].ID] = coursesWithoutMine[:10]
	}

	c.JSON(http.StatusOK, convertResult(recommended))
}

type Recommended struct {
	Course             Course
	RecommendedBecause []string
}

func convertResult(courses map[string][]Course) map[string]Recommended {
	result := make(map[string]Recommended)
	for id, recommendedCourses := range courses {
		for i := range recommendedCourses {
			if _, ok := result[recommendedCourses[i].ID]; ok {
				recbec := result[recommendedCourses[i].ID].RecommendedBecause
				res := Recommended{Course: recommendedCourses[i], RecommendedBecause: recbec}
				result[recommendedCourses[i].ID] = res
			} else {
				result[recommendedCourses[i].ID] = Recommended{Course: recommendedCourses[i], RecommendedBecause: []string{id}}
			}
		}
	}
	return result
}

func (s *State) FindCoursesAccordingFilter(c *gin.Context, filter interface{}, coursesCollection *mongo.Collection) ([]Course, error) {
	data, err := coursesCollection.Find(c, filter, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to find: %v", err)
	}

	var result []Course
	for data.Next(c) {
		l := Course{}
		err = data.Decode(&l)
		if err != nil {
			return nil, fmt.Errorf("unable to decode: %v", err)
		}
		result = append(result, l)
	}
	return result, nil
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