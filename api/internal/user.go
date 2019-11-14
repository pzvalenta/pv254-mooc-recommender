package internal

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID         *primitive.ObjectID `bson:"_id"`
	EnrolledIn []string            `bson:"enrolledIn"`
	Name       string              `bson:"name"`
}
