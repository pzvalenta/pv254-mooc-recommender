package internal

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

type State struct {
	DB         *mongo.Database
	customerID string
}

func NewState(customerID string) (*State, error) {
	DB, err := NewDatabase("localhost", "27017")
	if err != nil {
		return nil, fmt.Errorf("error creating state: %v", err)
	}
	return &State{DB: DB, customerID: customerID}, nil
}
