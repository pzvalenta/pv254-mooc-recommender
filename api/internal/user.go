package internal

import "go.mongodb.org/mongo-driver/bson/primitive"

//User ...
type User struct {
	ID         *primitive.ObjectID `bson:"_id"`
	EnrolledIn []string            `bson:"enrolledIn"`
	Name       string              `bson:"name"`
	Rating     []int64             `bson:"rating"`
	AuthID     string              `bson:"auth_id"`
}
