package internal

import (
	"context"
	"fmt"
	"log"
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
		{"$sample": bson.M{"size": 1}},
		{"$match": bson.M{"_id": bson.M{"$nin": myCourseIds}}}, //_id :{ $nin : [...] }
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
		fmt.Println(myCourses[i].Details.Language)
		filter := bson.D{
			{Key: "_id", Value: bson.D{{Key: "$nin", Value: myCourseIds}}},
			{Key: "subject", Value: myCourses[i].Subject},
			{Key: "categories", Value: bson.D{{Key: "$nin", Value: myCourses[i].Categories}}},
			//{Key: "details", Value: bson.D{{Key: "language", Value: myCourses[i].Details.Language}}},
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
	c.JSON(http.StatusOK, sorted[:Min(10, len(sorted))])
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

func (s *State) CategoryRecommending(c *gin.Context) {
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
