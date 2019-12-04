package internal

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

//State ...
type State struct {
	DB         *mongo.Database
}

//NewState ...
func NewState() (*State, error) {
	DB, err := NewDatabase("localhost", "27017")
	if err != nil {
		return nil, fmt.Errorf("error creating state: %v", err)
	}
	return &State{DB: DB}, nil
}
