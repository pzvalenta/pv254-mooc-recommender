package internal

import (
	"context"
	"log"
	"testing"
)

func TestSomething(t *testing.T) {

	s, err := NewState("5dceb44288861f034fc60b16")

	course1, _ := s.GetCourseByID("machine-learning-835")

	course2, _ := s.GetCourseByID("udacity-intro-to-machine-learning-2996")

	res := course1.tfidf([]Course{course2})
	log.Println(res)

	if err != nil {
		panic(err)
	}
}
func TestTfidf(t *testing.T) {
	s, err := NewState("5dceb44288861f034fc60b16")
	coursesCollection := s.DB.Collection("courses")
	var result []Course

	data, err := coursesCollection.Find(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	for data.Next(context.Background()) {
		l := Course{}
		err = data.Decode(&l)
		result = append(result, l)
	}
	log.Println("shiet")

}
