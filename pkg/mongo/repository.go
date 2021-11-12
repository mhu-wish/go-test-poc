package repository

import "go.mongodb.org/mongo-driver/mongo"

// Repository is the interface to data storage
type Repository interface {
	// Create new entry
	Create() error

	// Get all matched documents by criteria
	FindOne(filter interface{}) error

	// Find all desired documents, returns a cursor through which the caller can iterate and decode
	Find(filter interface{}) (*mongo.Cursor, error)
}
