package internal
//Recommended ...
type Recommended struct {
	Course             Course       `json:"course"`
	RecommendedBecause []Similarity `json:"recommendedBecause"`
}
//Similarity ...
type Similarity struct {
	CourseID   string
	Similarity float64
}
//OveralSimilarity ...
func (r *Recommended) OveralSimilarity() float64 {
	var overalSimilarity float64
	for _, similarity := range r.RecommendedBecause {
		overalSimilarity += similarity.Similarity
	}
	return overalSimilarity
}
//RecommendedAsArrayItem ...
type RecommendedAsArrayItem struct {
	CourseID          string      `json:"courseID"`
	Recommended       Recommended `json:"recommended"`
	OverallSimilarity float64     `json:"overallSimilarity"`
}

func fromMapWithSimilar(courses map[string][]SimilarCourse) map[string]*Recommended {
	result := make(map[string]*Recommended)
	for id, recommendedCourses := range courses {
		for i := range recommendedCourses {
			if _, ok := result[recommendedCourses[i].Course.ID]; ok {
				recbec := append(
					result[recommendedCourses[i].Course.ID].RecommendedBecause,
					Similarity{CourseID: id, Similarity: recommendedCourses[i].Similarity},
				)
				result[recommendedCourses[i].Course.ID].RecommendedBecause = recbec
			} else {
				result[recommendedCourses[i].Course.ID] =
					&Recommended{
						Course:             recommendedCourses[i].Course,
						RecommendedBecause: []Similarity{{CourseID: id, Similarity: recommendedCourses[i].Similarity}}}
			}
		}
	}
	return result
}
//FromRecommenedToSortedRecommended ...
func FromRecommenedToSortedRecommended(result map[string]*Recommended) []RecommendedAsArrayItem {
	var tmp []RecommendedAsArrayItem
	for recommendID, absolved := range result {
		sr := RecommendedAsArrayItem{
			CourseID:          recommendID,
			Recommended:       *absolved,
			OverallSimilarity: result[recommendID].OveralSimilarity(),
		}
		tmp = append(tmp, sr)
	}
	return tmp
}
//SortedByOverallSimilarity ...
type SortedByOverallSimilarity struct {
	sr []RecommendedAsArrayItem
}

func (s SortedByOverallSimilarity) Len() int      { return len(s.sr) }
func (s SortedByOverallSimilarity) Swap(i, j int) { s.sr[i], s.sr[j] = s.sr[j], s.sr[i] }
func (s SortedByOverallSimilarity) Less(i, j int) bool {
	return s.sr[i].OverallSimilarity > s.sr[j].OverallSimilarity
}
