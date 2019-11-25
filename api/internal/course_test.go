package internal

import (
	"context"
	"log"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestSomething(t *testing.T) {

	s, _ := NewState("5dceb44288861f034fc60b16")
	coursesCollection := s.DB.Collection("courses")
	c := context.Background()
	course1, _ := s.GetCourseByID("machine-learning-835")

	course2, _ := s.GetCourseByID("udacity-intro-to-machine-learning-2996")
	var courses []Course
	courses = append(courses, course1)
	courses = append(courses, course2)

	query := []bson.M{bson.M{"$sample": bson.M{"size": 20}}}

	dbCourses, _ := coursesCollection.Aggregate(c, query)
	dbCourses.All(c, &courses)

	res := course1.FindSimilar(courses, 0.0)

	log.Println(res)

}
