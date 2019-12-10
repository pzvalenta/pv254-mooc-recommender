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
		c.JSON(http.StatusBadRequest, fmt.Errorf("error creating id from hex: %v", err))
	}
	filter := bson.M{"_id": id}
	err = users.FindOne(c, filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Errorf("unable to find user's user IDs: %v", err))
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
		c.JSON(http.StatusBadRequest, fmt.Errorf("something went wrong: %v", err))
		return
	}

	c.JSON(http.StatusOK, u)

}

//EnrollUser ...
func (s *State) EnrollUser(c *gin.Context) {

	userID := c.Param("userId")
	courseID := c.Param("courseId")
	users := s.DB.Collection("users")

	id, err := primitive.ObjectIDFromHex(userID)
	var user User
	var course Course
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Errorf("error creating id from hex: %v", err))
		return
	}

	filter := bson.M{"_id": id}
	err = users.FindOne(c, filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Errorf("unable to find user's course IDs: %v", err))
		return
	}

	course, err = s.GetCourseByID(courseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Errorf("unable to find course: %v", err))
		return
	}
	isIn := false
	for _, b := range user.EnrolledIn {
		if b == course.ID {
			isIn = true
		}
	}
	if isIn {
		c.JSON(http.StatusBadRequest, "you are already enrolled in the course")
		return

	}
	user.EnrolledIn = append(user.EnrolledIn, course.ID)
	user.Rating = append(user.Rating, 0)

	update := bson.M{"$set": bson.M{"enrolledIn": user.EnrolledIn, "rating": user.Rating}}
	_, err = users.UpdateOne(c, filter, update)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Errorf("unable to update course: %v", err))
		return
	}
	c.JSON(http.StatusOK, user)
}

//RemoveUserEnrollment ...
func (s *State) RemoveUserEnrollment(c *gin.Context) {

	userID := c.Param("userId")
	courseID := c.Param("courseId")
	users := s.DB.Collection("users")

	id, err := primitive.ObjectIDFromHex(userID)
	var user User
	var course Course
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Errorf("error creating id from hex: %v", err))
		return
	}

	filter := bson.M{"_id": id}
	err = users.FindOne(c, filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Errorf("unable to find user's course IDs: %v", err))
		return
	}

	course, err = s.GetCourseByID(courseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Errorf("unable to find course: %v", err))
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
		c.JSON(http.StatusBadRequest, fmt.Errorf("unable to update course: %v", err))
		return
	}
	c.JSON(http.StatusOK, user)
}