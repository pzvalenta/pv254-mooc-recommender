package internal

type Recommended struct {
	Course             Course
	RecommendedBecause []Similarity
}

type Similarity struct {
	CourseID   string
	Similarity float64
}

func fromMapWithSimilar(courses map[string][]SimilarCourse) map[string]Recommended {
	result := make(map[string]Recommended)
	for id, recommendedCourses := range courses {
		for i := range recommendedCourses {
			if _, ok := result[recommendedCourses[i].Course.ID]; ok {
				recbec := append(result[recommendedCourses[i].Course.ID].RecommendedBecause, Similarity{CourseID: id, Similarity: recommendedCourses[i].Similarity})
				res := Recommended{Course: recommendedCourses[i].Course, RecommendedBecause: recbec}
				result[recommendedCourses[i].Course.ID] = res
			} else {
				result[recommendedCourses[i].Course.ID] =
					Recommended{
						Course:             recommendedCourses[i].Course,
						RecommendedBecause: []Similarity{{CourseID: id, Similarity: recommendedCourses[i].Similarity}}}
			}
		}
	}
	return result
}
