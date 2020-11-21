package store

import (
	"errors"
	"gobunny/store/model"
)

var (
	// ErrKeyAlreadyExists indicates a key has already been written to
	ErrKeyAlreadyExists = errors.New("key already exists")

	// ErrNotFound indicates a value was not found with the given key
	ErrNotFound = errors.New("key not found")
)

type (
	// Store is a CRUD-like interface for data structures backed by a data store
	Store interface {
		Creator
		Deleter
		Getter
		Updater
	}

	// ImmutableStore is a read-only interface for data structures backed by a data store
	ImmutableStore interface {
		Getter
	}

	// Creator creates a new data structure entry in the backing store
	Creator interface {
		Create(model.Key, model.Model) error
	}

	// Deleter deletes a data structure from the backing store
	Deleter interface {
		Delete(model.Key) error
	}

	// Getter fetches a data structure from the backing store
	Getter interface {
		Exists(model.Key) (bool, error)
		Get(model.Key, model.Model) error
	}

	// Updater updates an existing data structure in the backing store
	Updater interface {
		Update(model.Key, model.Model) error
	}
)
