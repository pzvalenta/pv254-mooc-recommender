package internal

import "go.mongodb.org/mongo-driver/bson/primitive"

//User ...
type User struct {
	ID         *primitive.ObjectID `json:"id" bson:"_id"`
	EnrolledIn []string            `json:"enrolledIn" bson:"enrolledIn"`
	Name       string              `json:"name" bson:"name"`
	Rating     []int64             `json:"rating" bson:"rating"`
	AuthID     string              `json:"authId" bson:"auth_id"`
}

//Review ...
type Review struct {
	ID       *primitive.ObjectID `json:"id" bson:"_id"`
	Text     string              `json:"text" bson:"text"`
	Rating   int64               `json:"rating" bson:"rating"`
	UserID   *primitive.ObjectID `json:"userId" bson:"user_id"`
	CourseID string              `json:"courseId" bson:"course_id"`
}

//Review ...
type ReviewModel struct {
	ID       *primitive.ObjectID `json:"id" bson:"_id"`
	Text     string              `json:"text" bson:"text"`
	Rating   int64               `json:"rating" bson:"rating"`
	User     User                `json:"user"`
	CourseID string              `json:"courseId" bson:"course_id"`
}
