package internal

import (
	"strings"

	"github.com/gin-gonic/gin"
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
	//StartDate        []string `json:"start date" bson:"start date"`
}

func ExtractSubjects(c *gin.Context, courses []Course) []string {
	subjects := make(map[string]interface{})
	for i := range courses {
		subjects[courses[i].Subject] = nil
	}

	var res []string
	for k := range subjects {
		res = append(res, k)
	}

	return res
}

func FindSimilar(c *gin.Context, courses, otherCourses []Course) []Course {
	var result []Course
	for i := range courses {
		for j := range otherCourses {
			if courses[i].isSimilar(&otherCourses[j]) > 0.7 {
				result = append(result, otherCourses[j])
			}
		}
	}
	return result
}

func (c *Course) isSimilar(c1 *Course) float64 {
	if c.ID == c1.ID {
		return 1.0
	}

	var result float64

	if c.Subject == c1.Subject {
		result += 1.0
	}

	result += float64(((len(c.Categories) / 100) * len(intersection(c.Categories, c1.Categories))) / 5)
	result += float64(((len(c.Schools) / 100) * len(intersection(c.Schools, c1.Schools))) / 5)
	result += float64(((len(strings.Split(c.Name, " ")) / 100) * len(intersection(strings.Split(c.Name, " "), strings.Split(c1.Name, " ")))) / 5)
	result += float64(((len(strings.Split(c.Description, " ")) / 100) * len(intersection(strings.Split(c.Description, " "), strings.Split(c1.Description, " ")))) / 5)

	return result
}

func intersection(a, b []string) []string {
	// interacting on the smallest list first can potentially be faster...but not by much, worse case is the same
	low, high := a, b
	if len(a) > len(b) {
		low = b
		high = a
	}

	var inter []string
	done := false
	for i, l := range low {
		for j, h := range high {
			// get future index values
			f1 := i + 1
			f2 := j + 1
			if l == h {
				inter = append(inter, h)
				if f1 < len(low) && f2 < len(high) {
					// if the future values aren't the same then that's the end of the intersection
					if low[f1] != high[f2] {
						done = true
					}
				}
				// we don't want to interate on the entire list everytime, so remove the parts we already looped on will make it faster each pass
				high = high[:j+copy(high[j:], high[j+1:])]
				break
			}
		}
		// nothing in the future so we are done
		if done {
			break
		}
	}
	return inter
}
