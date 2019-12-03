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
	myCourseIds, err := s.getMyCoursesIds(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "no content")
		return
	}

	query := []bson.M{
		bson.M{"$sample": bson.M{"size": 1}},
		bson.M{"$match": bson.M{"_id": bson.M{"$nin": myCourseIds}}}, //_id :{ $nin : [...] }
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
	myCourseIds, err := s.getMyCoursesIds(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "no content")
		return
	}

	myCourses, err := s.getMyCourses(c)
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
		}

		coursesFromOtherSubtree, err := s.findCoursesAccordingFilter(c, filter, coursesCollection)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "no content")
			return
		}

		similar := myCourses[i].FindSimilar(coursesFromOtherSubtree, 0.08)
		sort.Sort(SortedBySimilarity{coursesWithSimilarity: similar, course: &myCourses[i]})
		recommended[myCourses[i].ID] = similar
	}

	sorted := FromRecommenedToSortedRecommended(fromMapWithSimilar(recommended))
	sort.Sort(SortedByOverallSimilarity{sr: sorted})
	c.JSON(http.StatusOK, sorted)
}

//OverfittingCourses ...
func (s *State) OverfittingCourses(c *gin.Context) {
	myCourseIds, err := s.getMyCoursesIds(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "no content")
		return
	}

	myCourses, err := s.getMyCourses(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "no content")
		return
	}

	recommended := make(map[string][]SimilarCourse)
	coursesCollection := s.DB.Collection("courses")
	for i := range myCourses {
		filter := bson.D{
			{Key: "_id", Value: bson.D{{Key: "$nin", Value: myCourseIds}}},
		}
		coursesWithoutMine, err := s.findCoursesAccordingFilter(c, filter, coursesCollection)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "no content")
			return
		}

		similar := myCourses[i].FindSimilar(coursesWithoutMine, 0.1)
		sort.Sort(SortedBySimilarity{course: &myCourses[i], coursesWithSimilarity: similar})
		recommended[myCourses[i].ID] = similar
	}

	sorted := FromRecommenedToSortedRecommended(fromMapWithSimilar(recommended))
	sort.Sort(SortedByOverallSimilarity{sr: sorted})
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

func (s *State) getMyCoursesIds(c *gin.Context) ([]string, error) {
	users := s.DB.Collection("users")

	id, err := primitive.ObjectIDFromHex(s.customerID)
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

func (s *State) getMyCourses(c *gin.Context) ([]Course, error) {
	myCourseIDs, err := s.getMyCoursesIds(c)
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

// only in english
func (s *State) getAllCourses(c *gin.Context) ([]Course, error) {
	var res []Course
	filter := bson.D{{Key: "details.language", Value: bson.D{{Key: "$eq", Value: "English"}}}}
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
		bson.M{"$match": query},
		bson.M{
			"$sort": bson.M{
				"interested_count": -1,
				"review_count":     -1,
			},
		},
		bson.M{
			"$skip": int(page) * pageSize,
		},
		bson.M{
			"$limit": pageSize,
		},
	}
	dbCourses, err := coursesCollection.Aggregate(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "error")
		return
	}
	dbCourses.All(c, &result)
	c.JSON(http.StatusOK, result)

}

//GetAllSubjects ...
func (s *State) GetAllSubjects(c *gin.Context) {
	query := []bson.M{
		bson.M{"$project": bson.M{"subjects": bson.M{"$split": []interface{}{"$subject", ", "}}}},
		bson.M{"$unwind": bson.M{"path": "$subjects", "includeArrayIndex": "string", "preserveNullAndEmptyArrays": true}},
		bson.M{"$group": bson.M{"_id": nil, "unique_subjects": bson.M{"$addToSet": "$subjects"}}},
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
		bson.M{"$project": bson.M{"categoriess": "$categories", "subject": "$subject"}},
		bson.M{"$unwind": bson.M{"path": "$categoriess", "includeArrayIndex": "string", "preserveNullAndEmptyArrays": true}},
		bson.M{"$group": bson.M{"_id": "$subject", "unique_categories": bson.M{"$addToSet": "$categoriess"}}},
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
// return map    [courseAttributeValue] = IDF of said value
func (s *State) getUniqueAttributes(c *gin.Context, name string) (map[string]float64, error) {
	names := name + "s"

	query := []bson.M{
		//bson.M{"$project": bson.M{names: bson.M{"$split": []interface{}{"$" + name, ", "}}}},
		bson.M{"$project": bson.M{names: "$" + name}},
		bson.M{"$unwind": bson.M{"path": "$" + names, "includeArrayIndex": "string", "preserveNullAndEmptyArrays": true}},
		bson.M{"$group": bson.M{"_id": "$" + names, "count": bson.M{"$sum": 1}}},
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

	//fmt.Println(total)
	for i := range result {
		ret[result[i].ID] = math.Log10(float64(total) / float64(result[i].Count))
		//fmt.Println(result[i].ID, result[i].Count, ret[result[i].ID]) T
		// TODO there is one result which hase no ID
	}

	return ret, nil
}

// returns  userProfileVector dot product ( IDFvector dot product courseVector )
//                             [attributeName][attrValue]=IDFvalue
func predictCourseUser(IDFvectors map[string]map[string]float64, profile userProfile, course Course) float64 {
	//weighted course vector = IDFvector dot product courseVector
	/*
		weightedCourseVector := make(map[string]float64)
		weightedCourseVector[course.Subject] = IDFvectors["subject"][course.Subject]
		weightedCourseVector[course.Provider] = IDFvectors["provider"][course.Provider]
		for i := range course.Categories {
			weightedCourseVector[course.Categories[i]] = IDFvectors["categories"][course.Categories[i]]
		}
	*/

	var ret float64
	ret += IDFvectors["subject"][course.Subject] * profile.Subject[course.Subject]
	ret += IDFvectors["provider"][course.Provider] * profile.Provider[course.Provider]
	/*for i := range course.Categories {
		ret += IDFvectors["categories"][course.Categories[i]] * profile.Categories[course.Subject]
	}*/
	return ret
}

type userProfile struct {
	Categories map[string]float64
	Provider   map[string]float64
	Subject    map[string]float64
}

func (s *State) getUserProfile(c *gin.Context) (userProfile, error) {
	var profile userProfile
	profile.Categories = make(map[string]float64)
	profile.Provider = make(map[string]float64)
	profile.Subject = make(map[string]float64)

	myCourses, err := s.getMyCourses(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "no content")
		return profile, fmt.Errorf("%v", err)
	}

	for i := range myCourses {
		numberOfAttributes := 1 + 1 + len(myCourses[i].Categories) // 1 subject, 1 provider, x categories
		normalizedOccurence := 1 / math.Sqrt(float64(numberOfAttributes))

		// TODO add user rating of taken courses, multiply by rating, negative or positive
		profile.Subject[myCourses[i].Subject] += normalizedOccurence
		profile.Provider[myCourses[i].Provider] += normalizedOccurence
		for j := range myCourses[i].Categories {
			profile.Categories[myCourses[i].Categories[j]] += normalizedOccurence
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
	//fmt.Println(profile)

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

	allCourses, err := s.getAllCourses(c)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, "no content")
		return
	}

	for i := range allCourses {
		predictionValue := predictCourseUser(IDFvectors, profile, allCourses[i])
		//if predictionValue > 0 {
		fmt.Println(predictionValue, ":", allCourses[i].ID)
		//}
	}

	c.JSON(http.StatusOK, "general model recommendation")
}
