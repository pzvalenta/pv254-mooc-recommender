package internal

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//NewDatabase ...
func NewDatabase(host, port string) (*mongo.Database, error) {
	if envHost, ok := os.LookupEnv("DB_HOST"); ok {
		host = envHost
	}

	clientOptions := options.Client().ApplyURI("mongodb://" + host + ":" + port)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %v", err)
	}
	err = client.Connect(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error connecting to client: %v", err)
	}

	return client.Database("mydb"), nil
}

//RandomCourse ...
func (s *State) RandomCourse(c *gin.Context) {
	user_id := c.DefaultQuery("user_id", "5dc5715c70a18970fe47de7c")
	myCourseIds, err := s.getMyCoursesIds(c, user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "no content")
		return
	}

	query := []bson.M{
		{"$sample": bson.M{"size": 1}},
		{"$match": bson.M{"_id": bson.M{"$nin": myCourseIds}}}, //_id :{ $nin : [...] }
		{"$match": bson.M{"details.language": bson.M{"$eq": "English"}}},
	}

	coll := s.DB.Collection("courses")

	data, err := coll.Aggregate(c, query)

	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, "no content")
		return
	}

	var result []Course
	for data.Next(c) {
		l := Course{}
		err = data.Decode(&l)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "no content")
			return
		}
		result = append(result, l)
	}

	c.JSON(http.StatusOK, result)
}

//TaxonomyCourses ...
func (s *State) TaxonomyCourses(c *gin.Context) {
	user_id := c.DefaultQuery("user_id", "5dc5715c70a18970fe47de7c")
	myCourseIds, err := s.getMyCoursesIds(c, user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "no content")
		return
	}

	myCourses, err := s.getMyCourses(c, user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "no content")
		return
	}

	recommended := make(map[string][]SimilarCourse)
	coursesCollection := s.DB.Collection("courses")
	for i := range myCourses {
		filter := bson.D{
			{Key: "_id", Value: bson.D{{Key: "$nin", Value: myCourseIds}}},
			{Key: "subject", Value: myCourses[i].Subject},
			{Key: "categories", Value: bson.D{{Key: "$nin", Value: myCourses[i].Categories}}},
			{Key: "details.language", Value: "English"},
		}

		coursesFromOtherSubtree, err := s.findCoursesAccordingFilter(c, filter, coursesCollection)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "no content")
			return
		}

		similar := myCourses[i].FindSimilar(coursesFromOtherSubtree, 0.5)
		sort.Sort(SortedBySimilarity{coursesWithSimilarity: similar, course: &myCourses[i]})
		recommended[myCourses[i].ID] = similar
	}

	sorted := FromRecommenedToSortedRecommended(fromMapWithSimilar(recommended))
	sort.Sort(SortedByOverallSimilarity{sr: sorted})
	c.JSON(http.StatusOK, sorted[:Min(10, len(sorted))])
}

//OverfittingCourses ...
func (s *State) OverfittingCourses(c *gin.Context) {
	user_id := c.DefaultQuery("user_id", "5dc5715c70a18970fe47de7c")
	myCourseIds, err := s.getMyCoursesIds(c, user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "no content")
		return
	}

	myCourses, err := s.getMyCourses(c, user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "no content")
		return
	}

	recommended := make(map[string][]SimilarCourse)
	coursesCollection := s.DB.Collection("courses")
	for i := range myCourses {
		filter := bson.D{
			{Key: "_id", Value: bson.D{{Key: "$nin", Value: myCourseIds}}},
			{Key: "details.language", Value: "English"},
			{Key: "overview", Value: bson.D{{Key: "$nin", Value: []interface{}{nil, "", " ", "."}}}},
		}
		coursesWithoutMine, err := s.findCoursesAccordingFilter(c, filter, coursesCollection)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "no content")
			return
		}

		similar := myCourses[i].FindSimilar(coursesWithoutMine, 0.7)
		sort.Sort(SortedBySimilarity{course: &myCourses[i], coursesWithSimilarity: similar})
		recommended[myCourses[i].ID] = similar
	}

	sorted := FromRecommenedToSortedRecommended(fromMapWithSimilar(recommended))
	sort.Sort(SortedByOverallSimilarity{sr: sorted})
	c.JSON(http.StatusOK, sorted[:Min(10, len(sorted))])
}

// CategoryRecommending ...
func (s *State) CategoryRecommending(c *gin.Context) {
	user_id := c.DefaultQuery("user_id", "5dc5715c70a18970fe47de7c")
	myCourseIds, err := s.getMyCoursesIds(c, user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "no content")
		return
	}

	myCourses, err := s.getMyCourses(c, user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "no content")
		return
	}
	recommended := make(map[string][]PopularCourse)
	coursesCollection := s.DB.Collection("courses")
	for i := range myCourses {
		filter := bson.D{
			{Key: "_id", Value: bson.D{{Key: "$nin", Value: myCourseIds}}},
			{Key: "categories", Value: bson.D{{Key: "$in", Value: myCourses[i].Categories}}},
		}
		coursesWithoutMine, err := s.findCoursesAccordingFilter(c, filter, coursesCollection)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "no content")
			return
		}

		//include interested count
		popular := myCourses[i].FindSimilarAndPopular(coursesWithoutMine, 0.1)
		sort.Sort(SortedByPopularity{course: &myCourses[i], coursesWithPopularity: popular})
		recommended[myCourses[i].ID] = popular
	}

	sorted := FromRecommenedPopularToSortedRecommendedSorted(fromMapWithPopular(recommended))
	sort.Sort(SortedByOverallPopularity{sr: sorted})
	c.JSON(http.StatusOK, sorted[:Min(10, len(sorted))])
}

func (s *State) findCoursesAccordingFilter(c *gin.Context, filter interface{}, coursesCollection *mongo.Collection) ([]Course, error) {
	data, err := coursesCollection.Find(c, filter, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to find: %v", err)
	}

	var result []Course
	for data.Next(c) {
		l := Course{}
		err = data.Decode(&l)
		if err != nil {
			return nil, fmt.Errorf("unable to decode: %v", err)
		}
		result = append(result, l)
	}
	return result, nil
}

func (s *State) getMyCoursesIds(c *gin.Context, user_id string) ([]string, error) {
	users := s.DB.Collection("users")

	id, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		return nil, fmt.Errorf("error creating id from hex: %v", err)
	}

	filter := bson.M{"_id": id}
	data, err := users.Find(c, filter)
	if err != nil {
		return nil, fmt.Errorf("unable to find user's course IDs: %v", err)
	}

	l := User{}
	if data.Next(c) {
		err = data.Decode(&l)
		if err != nil {
			return nil, fmt.Errorf("unable to decode user's enrolledIn: %v", err)
		}
	}

	return l.EnrolledIn, nil
}

//GetCourseByID ...
func (s *State) GetCourseByID(id string) (Course, error) {
	coursesCollection := s.DB.Collection("courses")
	var result Course

	filter := bson.M{"_id": id}
	err := coursesCollection.FindOne(context.Background(), filter).Decode(&result)
	return result, err
}

func (s *State) getMyCourses(c *gin.Context, user_id string) ([]Course, error) {
	myCourseIDs, err := s.getMyCoursesIds(c, user_id)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	ids := bson.A{}
	for _, id := range myCourseIDs {
		ids = append(ids, id)
	}

	var res []Course
	filter := bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: ids}}}}
	courseCollection := s.DB.Collection("courses")
	data, err := courseCollection.Find(c, filter, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to find/decode courses %v", err)
	}

	for data.Next(c) {
		l := Course{}
		err = data.Decode(&l)
		if err != nil {
			return nil, fmt.Errorf("unable to decode user's enrolledIn: %v", err)
		}
		res = append(res, l)
	}

	return res, nil
}

func (s *State) getIdf() (*map[string]float64, error) {
	idfCol := s.DB.Collection("idf")
	c := context.Background()
	data, err := idfCol.Find(c, bson.M{})
	res := make(map[string]float64)
	if err != nil {
		return nil, fmt.Errorf("unable to find/decode idfs %v", err)
	}

	for data.Next(c) {
		l := WordIdf{}
		err = data.Decode(&l)
		if err != nil {
			return nil, fmt.Errorf("unable to decode idf: %v", err)
		}
		res[l.Word] = l.Value
	}

	return &res, nil
}

func (s *State) getAllCoursesWithoutMine(c *gin.Context, user_id string) ([]Course, error) {
	myCourseIDs, err := s.getMyCoursesIds(c, user_id)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	ids := bson.A{}
	for _, id := range myCourseIDs {
		ids = append(ids, id)
	}

	//only in english
	var res []Course
	filter := bson.D{{Key: "_id", Value: bson.D{{Key: "$nin", Value: ids}}},
		{Key: "details.language", Value: bson.D{{Key: "$eq", Value: "English"}}}}

	courseCollection := s.DB.Collection("courses")
	data, err := courseCollection.Find(c, filter, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to find/decode courses %v", err)
	}

	for data.Next(c) {
		l := Course{}
		err = data.Decode(&l)
		if err != nil {
			return nil, fmt.Errorf("unable to decode user's enrolledIn: %v", err)
		}
		res = append(res, l)
	}

	return res, nil
}

func (s *State) getMyRatings(c *gin.Context, user_id string) (map[string]int64, error) {
	users := s.DB.Collection("users")

	id, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		return nil, fmt.Errorf("error creating id from hex: %v", err)
	}

	filter := bson.M{"_id": id}
	data, err := users.Find(c, filter)
	if err != nil {
		return nil, fmt.Errorf("unable to find user's course IDs: %v", err)
	}

	l := User{}
	if data.Next(c) {
		err = data.Decode(&l)
		if err != nil {
			return nil, fmt.Errorf("unable to decode user's enrolledIn: %v", err)
		}
	}

	ret := make(map[string]int64)

	for i, c := range l.EnrolledIn {
		ret[c] = l.Rating[i]
	}

	return ret, nil
}

// GetCoursebByID ...
func (s *State) GetCoursebByID(c *gin.Context) {
	id := c.Param("id")
	coursesCollection := s.DB.Collection("courses")
	var result Course

	filter := bson.D{
		{Key: "_id", Value: bson.D{{Key: "$eq", Value: id}}},
	}

	err := coursesCollection.FindOne(c, filter).Decode(&result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "no content")
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetCoursesByQuery ...
func (s *State) GetCoursesByQuery(c *gin.Context) {
	language := c.DefaultQuery("language", "English")
	pageString := c.DefaultQuery("page", "0")
	subject := c.DefaultQuery("subject", "")
	provider := c.DefaultQuery("provider", "")
	category := c.DefaultQuery("category", "")
	school := c.DefaultQuery("school", "")
	pageSize := 20

	page, err := strconv.ParseUint(pageString, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	query := bson.M{
		"details.language": language,
	}
	if subject != "" {
		query["subject"] = subject
	}
	if category != "" {
		query["categories"] = category
	}
	if provider != "" {
		query["provider"] = provider
	}
	if school != "" {
		query["schools"] = school
	}

	coursesCollection := s.DB.Collection("courses")
	var result []Course
	filter := []bson.M{
		{"$match": query},
		{
			"$sort": bson.M{
				"interested_count": -1,
				"review_count":     -1,
			},
		},
		{
			"$skip": int(page) * pageSize,
		},
		{
			"$limit": pageSize,
		},
	}
	dbCourses, err := coursesCollection.Aggregate(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "no content")
		return
	}
	err = dbCourses.All(c, &result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "no content")
		return
	}
	c.JSON(http.StatusOK, result)
}

//GetAllSubjects ...
func (s *State) GetAllSubjects(c *gin.Context) {
	query := []bson.M{
		{"$project": bson.M{"subjects": bson.M{"$split": []interface{}{"$subject", ", "}}}},
		{"$unwind": bson.M{"path": "$subjects", "includeArrayIndex": "string", "preserveNullAndEmptyArrays": true}},
		{"$group": bson.M{"_id": nil, "unique_subjects": bson.M{"$addToSet": "$subjects"}}},
	}

	coll := s.DB.Collection("courses")

	data, err := coll.Aggregate(c, query)

	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, "no content")
		return
	}

	type subjects struct {
		UniqueSubjects []string `json:"unique_subjects" bson:"unique_subjects"`
	}
	var result []subjects

	err = data.All(c, &result)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, "no content")
		return
	}
	c.JSON(http.StatusOK, result[0])
}

//GetAllCategories ...
func (s *State) GetAllCategories(c *gin.Context) {
	query := []bson.M{
		{"$project": bson.M{"categoriess": "$categories", "subject": "$subject"}},
		{"$unwind": bson.M{"path": "$categoriess", "includeArrayIndex": "string", "preserveNullAndEmptyArrays": true}},
		{"$group": bson.M{"_id": "$subject", "unique_categories": bson.M{"$addToSet": "$categoriess"}}},
	}

	coll := s.DB.Collection("courses")

	data, err := coll.Aggregate(c, query)

	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, "no content")
		return
	}

	type categories struct {
		ID               string   `json:"_id" bson:"_id"`
		UniqueCategories []string `json:"unique_categories" bson:"unique_categories"`
	}
	var result []categories

	err = data.All(c, &result)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, "no content")
		return
	}
	c.JSON(http.StatusOK, result)
}

type attribute struct {
	ID    string `json:"_id" bson:"_id"`
	Count int    `json:"count" bson:"count"`
	IDF   float64
}

// TODO refactor into two functions, one universal which returns array of attribute and second one which returns IDF
// it's basically the same as GetAllCategories and GetAllSubjects
// return map    [courseAttributeValue] = IDF of said value
func (s *State) getUniqueAttributes(c *gin.Context, name string) (map[string]float64, error) {
	names := name + "s"

	query := []bson.M{
		bson.M{"$project": bson.M{names: "$" + name}},
		bson.M{"$unwind": bson.M{"path": "$" + names, "includeArrayIndex": "string", "preserveNullAndEmptyArrays": true}},
		bson.M{"$group": bson.M{"_id": "$" + names, "count": bson.M{"$sum": 1}}},
		bson.M{"$match": bson.M{"_id": bson.M{"$ne": nil}}},
	}

	coll := s.DB.Collection("courses")
	data, err := coll.Aggregate(c, query)
	if err != nil {
		return nil, err
	}

	var result []attribute
	err = data.All(c, &result)
	if err != nil {
		return nil, err
	}

	total := 0
	for i := range result {
		total += result[i].Count
	}

	ret := make(map[string]float64, total)

	for i := range result {
		ret[result[i].ID] = math.Log10(float64(total) / float64(result[i].Count))
		// TODO there is one result which hase no ID
	}

	return ret, nil
}

// returns  userProfileVector dot product ( IDFvector dot product courseVector )
//                             [attributeName][attrValue]=IDFvalue
func predictCourseUser(IDFvectors map[string]map[string]float64, profile userProfile, course Course) float64 {
	// based on input course
	numberOfAttributes := 1 + len(course.Categories) + len(course.Schools) + len(course.Teachers)
	normalizedOccurence := 1 / math.Sqrt(float64(numberOfAttributes))

	var ret float64

	ret += IDFvectors["subject"][course.Subject] * profile.Subject[course.Subject] * normalizedOccurence
	for i := range course.Categories {
		ret += IDFvectors["categories"][course.Categories[i]] * profile.Categories[course.Categories[i]]
	}
	for i := range course.Schools {
		ret += IDFvectors["schools"][course.Schools[i]] * profile.Schools[course.Schools[i]] * 0.5 /// school should have lower weight
	}
	for i := range course.Teachers {
		ret += IDFvectors["teachers"][course.Teachers[i]] * profile.Teachers[course.Teachers[i]]
	}
	return ret
}

type userProfile struct {
	Categories map[string]float64
	Provider   map[string]float64
	Subject    map[string]float64
	Teachers   map[string]float64
	Schools    map[string]float64
}

func (s *State) getUserProfile(c *gin.Context) (userProfile, error) {
	var profile userProfile
	profile.Categories = make(map[string]float64)
	profile.Provider = make(map[string]float64)
	profile.Subject = make(map[string]float64)
	profile.Teachers = make(map[string]float64)
	profile.Schools = make(map[string]float64)

	user_id := c.DefaultQuery("user_id", "5dc5715c70a18970fe47de7c")

	myCourses, err := s.getMyCourses(c, user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "no content")
		return profile, fmt.Errorf("%v", err)
	}

	myRatings, err := s.getMyRatings(c, user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "no content")
		return profile, fmt.Errorf("%v", err)
	}

	for i := range myCourses {
		numberOfAttributes := 1 + len(myCourses[i].Categories) + len(myCourses[i].Schools) + len(myCourses[i].Teachers)
		normalizedOccurence := 1 / math.Sqrt(float64(numberOfAttributes))

		profile.Subject[myCourses[i].Subject] += normalizedOccurence * float64(myRatings[myCourses[i].ID])
		profile.Provider[myCourses[i].Provider] += normalizedOccurence * float64(myRatings[myCourses[i].ID])
		for j := range myCourses[i].Categories {
			profile.Categories[myCourses[i].Categories[j]] += normalizedOccurence * float64(myRatings[myCourses[i].ID])
		}
		for j := range myCourses[i].Schools {
			profile.Schools[myCourses[i].Schools[j]] += normalizedOccurence * float64(myRatings[myCourses[i].ID])
		}
		for j := range myCourses[i].Teachers {
			profile.Teachers[myCourses[i].Teachers[j]] += normalizedOccurence * float64(myRatings[myCourses[i].ID])
		}
	}
	return profile, nil
}

// GeneralModelCourses ...
// https://www.analyticsvidhya.com/blog/2015/08/beginners-guide-learn-content-based-recommender-systems/
func (s *State) GeneralModelCourses(c *gin.Context) {
	profile, err := s.getUserProfile(c)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, "no content")
		return
	}

	IDFvectors := make(map[string]map[string]float64)

	IDFvectors["subject"], err = s.getUniqueAttributes(c, "subject")
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, "no content")
		return
	}
	IDFvectors["provider"], err = s.getUniqueAttributes(c, "provider")
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, "no content")
		return
	}
	IDFvectors["categories"], err = s.getUniqueAttributes(c, "categories")
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, "no content")
		return
	}
	IDFvectors["schools"], err = s.getUniqueAttributes(c, "schools")
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, "no content")
		return
	}
	IDFvectors["teachers"], err = s.getUniqueAttributes(c, "teachers")
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, "no content")
		return
	}

	user_id := c.DefaultQuery("user_id", "5dc5715c70a18970fe47de7c")
	allCourses, err := s.getAllCoursesWithoutMine(c, user_id)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, "no content")
		return
	}

	var tmp []RecommendedAsArrayItem

	for i := range allCourses {
		predictionValue := predictCourseUser(IDFvectors, profile, allCourses[i])

		if predictionValue > 0 {
			rec := RecommendedSimilar{
				Course: allCourses[i],
				//RecommendedBecause: []Similarity `json:"recommendedBecause"`
			}

			sr := RecommendedAsArrayItem{
				CourseID:          allCourses[i].ID,
				Recommended:       rec,
				OverallSimilarity: predictionValue,
			}

			tmp = append(tmp, sr)
		}
	}

	sort.Sort(SortedByOverallSimilarity{sr: tmp})

	c.JSON(http.StatusOK, tmp[:Min(10, len(tmp))])
}
