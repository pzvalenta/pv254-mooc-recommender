package internal

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//GetUserByID ...
func (s *State) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")

	users := s.DB.Collection("users")

	var user User

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error creating id from hex: %v", err)})
	}
	filter := bson.M{"_id": id}
	err = users.FindOne(c, filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("unable to find user's user IDs: %v", err)})
		return
	}
	c.JSON(http.StatusOK, user)
}

//GetUserByAuthID ...
func (s *State) GetUserByAuthID(c *gin.Context) {
	id := c.Param("authId")

	users := s.DB.Collection("users")

	var user User

	filter := bson.M{"auth_id": id}
	err := users.FindOne(c, filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("unable to find user's user IDs: %v", err)})
		return
	}
	c.JSON(http.StatusOK, user)
}

//CreateUser ...
func (s *State) CreateUser(c *gin.Context) {
	name := c.Param("name")
	authID := c.Param("authId")
	users := s.DB.Collection("users")
	u := User{Name: name, AuthID: authID, Rating: []int64{}, EnrolledIn: []string{}}
	res, err := users.InsertOne(c, bson.M{"name": name, "auth_id": authID, "rating": []int64{}, "enrolledIn": []string{}})

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		u.ID = &oid
	}
	fmt.Println(res)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("something went wrong: %v", err)})
		return
	}

	c.JSON(http.StatusOK, u)

}

//EnrollUser ...
func (s *State) EnrollUser(c *gin.Context) {

	authID := c.Param("authId")
	courseID := c.Param("courseId")
	users := s.DB.Collection("users")
	var user User
	var course Course

	filter := bson.M{"auth_id": authID}
	err := users.FindOne(c, filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("unable to find user by ID: %v", err)})
		return
	}

	course, err = s.GetCourseByID(courseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("unable to find course: %v", err)})
		return
	}
	isIn := false
	for _, b := range user.EnrolledIn {
		if b == course.ID {
			isIn = true
		}
	}
	if isIn {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you are already enrolled in the course"})
		return

	}
	user.EnrolledIn = append(user.EnrolledIn, course.ID)
	user.Rating = append(user.Rating, 0)

	update := bson.M{"$set": bson.M{"enrolledIn": user.EnrolledIn, "rating": user.Rating}}
	_, err = users.UpdateOne(c, filter, update)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("unable to update course: %v", err)})
		return
	}
	c.JSON(http.StatusOK, user)
}

//RemoveUserEnrollment ...
func (s *State) RemoveUserEnrollment(c *gin.Context) {

	authID := c.Param("authId")
	courseID := c.Param("courseId")
	users := s.DB.Collection("users")

	var user User
	var course Course

	filter := bson.M{"auth_id": authID}
	err := users.FindOne(c, filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("unable to find user's course IDs: %v", err)})
		return
	}

	course, err = s.GetCourseByID(courseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("unable to find course: %v", err)})
		return
	}
	index := -1
	newEnrollment := make([]string, 0)
	newRatings := make([]int64, 0)

	for i, b := range user.EnrolledIn {
		if b == course.ID {
			index = i
		} else {
			newEnrollment = append(newEnrollment, b)
			newRatings = append(newRatings, user.Rating[i])
		}
	}
	if index == -1 {
		c.JSON(http.StatusBadRequest, "you are not enrolled in the course")
		return
	}

	user.EnrolledIn = newEnrollment
	user.Rating = newRatings

	update := bson.M{"$set": bson.M{"enrolledIn": user.EnrolledIn, "rating": user.Rating}}
	_, err = users.UpdateOne(c, filter, update)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("unable to update course: %v", err)})
		return
	}
	c.JSON(http.StatusOK, user)
}

//GetUserCoursesByAuth this method needs rework
func (s *State) GetUserCoursesByAuth(c *gin.Context) {
	authID := c.Param("authId")
	users := s.DB.Collection("users")

	var user User
	filter := bson.M{"auth_id": authID}
	err := users.FindOne(c, filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("unable to find user's course IDs: %v", err)})
		return
	}
	courses, err := s.getMyCourses(c, user.ID.Hex())

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("something went wrong with courses: %v", err)})
		return
	}
	c.JSON(http.StatusOK, courses)

}

//GetUserCourses ...
func (s *State) GetUserCourses(c *gin.Context) {
	userID := c.Param("id")
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error creating id from hex: %v", err)})
		return
	}
	courses, err := s.getMyCourses(c, id.Hex())

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("something went wrong with courses: %v", err)})
		return
	}
	c.JSON(http.StatusOK, courses)

}

type reviewForm struct {
	Text   string `json:"text"`
	Rating int64  `json:"rating"`
}

//PostReview ...
func (s *State) PostReview(c *gin.Context) {
	authID := c.Param("authId")
	courseID := c.Param("courseId")

	users := s.DB.Collection("users")
	reviews := s.DB.Collection("reviews")
	var user User
	var course Course
	var r reviewForm
	c.BindJSON(&r)

	filter := bson.M{"auth_id": authID}
	err := users.FindOne(c, filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("unable to find user's course IDs: %v", err)})
		return
	}
	course, err = s.GetCourseByID(courseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("unable to find course: %v", err)})
		return
	}

	newReview := Review{UserID: user.ID, CourseID: course.ID, Text: r.Text, Rating: r.Rating}
	index := -1

	for i, b := range user.EnrolledIn {
		if b == course.ID {
			user.Rating[i] = newReview.Rating
			index = i
		}
	}
	if index == -1 {
		c.JSON(http.StatusBadRequest, "you are not enrolled in the course")
		return
	}
	res, err := reviews.InsertOne(c, bson.M{"text": newReview.Text, "user_id": newReview.UserID, "rating": newReview.Rating, "course_id": newReview.CourseID})

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		newReview.ID = &oid
	}
	update := bson.M{"$set": bson.M{"rating": user.Rating}}
	_, err = users.UpdateOne(c, filter, update)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("unable to update course: %v", err)})
		return
	}
	c.JSON(http.StatusOK, newReview)

}

//GetCourseReviews ...
func (s *State) GetCourseReviews(c *gin.Context) {
	courseID := c.Param("courseId")
	users := s.DB.Collection("users")

	reviews := s.DB.Collection("reviews")

	filter := bson.M{"course_id": courseID}
	data, err := reviews.Find(c, filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("unable to find course: %v", err)})
		return
	}

	var revs []Review
	for data.Next(c) {
		var r Review
		err = data.Decode(&r)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("unable to decode: %v", err)})
			return
		}
		revs = append(revs, r)
	}
	var revModel []ReviewModel
	userRev := make(map[string]User)
	var user User
	for _, rev := range revs {
		if val, ok := userRev[rev.UserID.Hex()]; ok {
			user = val
		} else {
			filter := bson.M{"_id": rev.UserID}
			err := users.FindOne(c, filter).Decode(&user)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("unable to find user's course IDs: %v", err)})
				return
			}
		}
		revModel = append(revModel, ReviewModel{CourseID: rev.CourseID,
			ID: rev.ID, Text: rev.Text, Rating: rev.Rating, User: user})
	}

	c.JSON(http.StatusOK, revModel)

}

//DeleteReview ...
func (s *State) DeleteReview(c *gin.Context) {
	authID := c.Param("authId")
	reviewID := c.Param("reviewId")

	users := s.DB.Collection("users")
	reviews := s.DB.Collection("reviews")
	var user User
	var review Review

	id, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error creating id from hex: %v", err)})
	}
	filter := bson.M{"_id": id}
	err = reviews.FindOne(c, filter).Decode(&review)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("unable to find user's user IDs: %v", err)})
		return
	}

	filter = bson.M{"auth_id": authID}
	err = users.FindOne(c, filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("unable to find user's course IDs: %v", err)})
		return
	}

	if review.UserID.Hex() != user.ID.Hex() {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("thats not your review: %v", err)})
		return
	}
	update := bson.M{"$set": bson.M{"rating": user.Rating}}
	_, err = users.UpdateOne(c, filter, update)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("unable to update course: %v", err)})
		return
	}
	filter = bson.M{"_id": id}

	_, err = reviews.DeleteOne(c, filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("unable to delete : %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{})

}
