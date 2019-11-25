package internal

import (
	"math"
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

//SortedBySimilarity ...
type SortedBySimilarity struct {
	coursesWithSimilarity []SimilarCourse
	course                *Course
}

//SimilarCourse ...
type SimilarCourse struct {
	Course     Course
	Similarity float64
}

func (s SortedBySimilarity) Len() int { return len(s.coursesWithSimilarity) }
func (s SortedBySimilarity) Swap(i, j int) {
	s.coursesWithSimilarity[i], s.coursesWithSimilarity[j] = s.coursesWithSimilarity[j], s.coursesWithSimilarity[i]
}
func (s SortedBySimilarity) Less(i, j int) bool {
	return s.coursesWithSimilarity[i].Similarity < s.coursesWithSimilarity[j].Similarity
}

//FindSimilar ...
func (c *Course) FindSimilar(courses []Course, similarityThresold float64) []SimilarCourse {
	var result []SimilarCourse

	var courseOverviews []string
	for i := range courses {
		courseOverviews = append(courseOverviews, courses[i].Overview)
	}
	idf := computeIdf(courseOverviews)
	var simVals []float64

	maxVal := 0.0
	for i := range courses {
		simVal := c.isSimilar(&courses[i], idf)
		maxVal = math.Max(simVal, maxVal)
		simVals = append(simVals, simVal)
	}
	for i := range simVals {
		simVals[i] = math.Abs((simVals[i] / maxVal) - 1)
	}
	for i := range simVals {
		simVal := simVals[i]
		if simVal > similarityThresold {
			result = append(result, SimilarCourse{Course: courses[i], Similarity: simVal})
		}
	}
	return result
}

func (c *Course) tfidf(idf map[string]float64) map[string]float64 {
	var tfidf map[string]float64
	tf := computeTf(c.Overview)
	for word, val := range tf {
		tfidf[word] = val * idf[word]
	}
	return tfidf
}

func (c *Course) isSimilar(c1 *Course, idf map[string]float64) float64 {
	if c.ID == c1.ID {
		return 1.0
	}
	tfidf1 := c.tfidf(idf)
	tfidf2 := c1.tfidf(idf)

	res := 0.0
	var wordList map[string]bool

	for word := range tfidf1 {
		wordList[word] = true
	}
	for word := range tfidf2 {
		wordList[word] = true
	}
	for word := range wordList {
		val1, val2 := 0.0, 0.0
		if num, ok := tfidf1[word]; ok {
			val1 = num
		}
		if num, ok := tfidf2[word]; ok {
			val2 = num
		}
		res += math.Pow(val1-val2, 2)
	}
	return math.Sqrt(res)
}
