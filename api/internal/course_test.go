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

	var courses []Course

	query := []bson.M{
		bson.M{
			"$match": bson.M{
				"details.language": "English",
				"overview":         bson.M{"$nin": []interface{}{nil, "", " ", "."}},
			},
		},
		// bson.M{"$sample": bson.M{"size": 10}},
	}

	dbCourses, err := coursesCollection.Aggregate(c, query)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	dbCourses.All(c, &courses)

	res := course1.FindSimilar(courses, 20)

	log.Println(res)

}

func TestIDF(t *testing.T) {
	x := []string{
		". It is often used as a weighting factor information in searches of information retrieval, text mining, and user modeling.",
		", is a numerical statistic that is intended to reflect how important a word is to a document in a collection or corpus",
		"In information retrieval, tf–idf or TFIDF, short for term frequency–inverse document frequency",
	}
	res := computeIdf(x)
	log.Println(res)

}

func TestCreateIdfList(t *testing.T) {
	s, _ := NewState("5dceb44288861f034fc60b16")
	coursesCollection := s.DB.Collection("courses")
	c := context.Background()
	var courses []Course
	query := []bson.M{
		bson.M{
			"$match": bson.M{
				"details.language": "English",
				"overview":         bson.M{"$nin": []interface{}{nil, "", " ", "."}},
			},
		},
	}

	dbCourses, err := coursesCollection.Aggregate(c, query)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	dbCourses.All(c, &courses)

	var overviews []string
	for c := range courses {
		overviews = append(overviews, courses[c].Overview)
	}

	idfCollection := s.DB.Collection("idf")
	idfCollection.Drop(c)
	resMap := computeIdf(overviews)
	var resStructs []interface{}
	for x, val := range resMap {
		resStructs = append(resStructs, WordIdf{x, val})
	}
	idfCollection.InsertMany(c, resStructs)

	log.Println("end")
}

type WordIdf struct {
	Word  string
	Value float64
}
