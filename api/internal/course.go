package internal

import (
	"strings"
)

// Course ...
type Course struct {
	ID              string   `json:"id" bson:"_id"`
	Categories      []string `json:"categories" bson:"categories"`
	Description     string   `json:"description" bson:"description"`
	Details         Details  `json:"details" bson:"details"`
	InterestedCount int32    `json:"interested_count" bson:"interested_count"`
	Link            string   `json:"link" bson:"link"`
	Name            string   `json:"name" bson:"name"`
	Overview        string   `json:"overview" bson:"overview"`
	Provider        string   `json:"provider" bson:"provider"`
	Rating          *float64 `json:"rating" bson:"rating"`
	ReviewCount     int32    `json:"review_count" bson:"review_count"`
	Schools         []string `json:"schools" bson:"schools"`
	Subject         string   `json:"subject" bson:"subject"`
	Syllabus        *string  `json:"syllabus" bson:"syllabus"`
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
	StartDate        []string `json:"start_date" bson:"start_date"`
}

type BySimilarity struct {
	coursesWithSimilarity []SimilarCourse
	course                *Course
}

func (s BySimilarity) Len() int { return len(s.coursesWithSimilarity) }
func (s BySimilarity) Swap(i, j int) {
	s.coursesWithSimilarity[i], s.coursesWithSimilarity[j] = s.coursesWithSimilarity[j], s.coursesWithSimilarity[i]
}
func (s BySimilarity) Less(i, j int) bool {
	return s.coursesWithSimilarity[i].Similarity < s.coursesWithSimilarity[j].Similarity
}

type SimilarCourse struct {
	Course     Course
	Similarity float64
}

func (c *Course) FindSimilar(courses []Course, similarityThresold float64) []SimilarCourse {
	var result []SimilarCourse
	for i := range courses {
		if c.isSimilar(&courses[i]) > similarityThresold {
			result = append(result, SimilarCourse{Course: courses[i], Similarity: c.isSimilar(&courses[i])})
		}
	}
	return result
}

func getStopWords() map[string]string {
	return map[string]string{" a ": " ", " and ": " ", " the ": " ", " of ": " ", " is ": " ", " are ": " ",
		" in ": " ", " to ": " ", " from ": " ", " on ": " ", ".": "", ":": ""}
}

func (c *Course) isSimilar(c1 *Course) float64 {
	if c.ID == c1.ID {
		return 1.0
	}

	numberOfAttributes := 1.0

	name1 := c.Name
	name2 := c1.Name
	for word, newWord := range getStopWords() {
		name1 = strings.Replace(name1, word, newWord, -1)
		name2 = strings.Replace(name2, word, newWord, -1)
	}

	result := float64(len(intersection(strings.Split(name1, " "), strings.Split(name2, " ")))) / float64(len(strings.Split(name1, " "))) / numberOfAttributes
	return result
}
