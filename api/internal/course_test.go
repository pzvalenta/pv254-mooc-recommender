package internal

import (
	"fmt"
	"log"
	"testing"
)

func TestSomething(t *testing.T) {

	s, err := NewState("5dceb44288861f034fc60b16")

	course1, _ := s.GetCourseByID("machine-learning-835")

	course2, _ := s.GetCourseByID("udacity-intro-to-machine-learning-2996")

	res := course1.tfidf(course2)
	log.Println(fmt.Sprintf("courses are similar with %f prob", res))

	if err != nil {
		panic(err)
	}
}
