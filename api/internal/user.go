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
