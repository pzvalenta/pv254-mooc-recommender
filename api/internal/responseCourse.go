package internal

//Recommended ...
type RecommendedSimilar struct {
	Course             Course       `json:"course"`
	RecommendedBecause []Similarity `json:"recommendedBecause"`
}

type RecommendedPopular struct {
	Course             Course       `json:"course"`
	RecommendedBecause []Popularity `json:"recommendedBecause"`
}

//Similarity ...
type Similarity struct {
	CourseID   string
	Similarity float64
}

type Popularity struct {
	CourseID   string
	Popularity float64
}

//OveralSimilarity ...
func (r *RecommendedSimilar) OveralSimilarity() float64 {
	var overalSimilarity float64
	for _, similarity := range r.RecommendedBecause {
		overalSimilarity += similarity.Similarity
	}
	return overalSimilarity
}

func (r *RecommendedPopular) OveralPopularity() float64 {
	var overalPopularity float64
	for _, similarity := range r.RecommendedBecause {
		overalPopularity += similarity.Popularity
	}
	return overalPopularity
}

//RecommendedAsArrayItem ...
type RecommendedAsArrayItem struct {
	CourseID          string             `json:"courseID"`
	Recommended       RecommendedSimilar `json:"recommended"`
	OverallSimilarity float64            `json:"overallSimilarity"`
}

type RecommendedAsArrayItemPopular struct {
	CourseID          string             `json:"courseID"`
	Recommended       RecommendedPopular `json:"recommended"`
	OverallPopularity float64            `json:"overallSimilarity"`
}

func fromMapWithSimilar(courses map[string][]SimilarCourse) map[string]*RecommendedSimilar {
	result := make(map[string]*RecommendedSimilar)
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
					&RecommendedSimilar{
						Course:             recommendedCourses[i].Course,
						RecommendedBecause: []Similarity{{CourseID: id, Similarity: recommendedCourses[i].Similarity}}}
			}
		}
	}
	return result
}

func fromMapWithPopular(courses map[string][]PopularCourse) map[string]*RecommendedPopular {
	result := make(map[string]*RecommendedPopular)
	for id, recommendedCourses := range courses {
		for i := range recommendedCourses {
			if _, ok := result[recommendedCourses[i].Course.ID]; ok {
				recbec := append(
					result[recommendedCourses[i].Course.ID].RecommendedBecause,
					Popularity{CourseID: id, Popularity: recommendedCourses[i].Popularity},
				)
				result[recommendedCourses[i].Course.ID].RecommendedBecause = recbec
			} else {
				result[recommendedCourses[i].Course.ID] =
					&RecommendedPopular{
						Course:             recommendedCourses[i].Course,
						RecommendedBecause: []Popularity{{CourseID: id, Popularity: recommendedCourses[i].Popularity}}}
			}
		}
	}
	return result
}

//FromRecommenedToSortedRecommended ...
func FromRecommenedToSortedRecommended(result map[string]*RecommendedSimilar) []RecommendedAsArrayItem {
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

func FromRecommenedPopularToSortedRecommendedSorted(result map[string]*RecommendedPopular) []RecommendedAsArrayItemPopular {
	var tmp []RecommendedAsArrayItemPopular
	for recommendID, absolved := range result {
		sr := RecommendedAsArrayItemPopular{
			CourseID:          recommendID,
			Recommended:       *absolved,
			OverallPopularity: result[recommendID].OveralPopularity(),
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

type SortedByOverallPopularity struct {
	sr []RecommendedAsArrayItemPopular
}

func (s SortedByOverallPopularity) Len() int      { return len(s.sr) }
func (s SortedByOverallPopularity) Swap(i, j int) { s.sr[i], s.sr[j] = s.sr[j], s.sr[i] }
func (s SortedByOverallPopularity) Less(i, j int) bool {
	return s.sr[i].OverallPopularity > s.sr[j].OverallPopularity
}
