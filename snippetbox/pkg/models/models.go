// Package models : This is intended to handle the models definitions to the whole application.
// must be included into other packages
package models

import (
	"errors"
	"time"
)

// ErrNoRecord - Custom error when there a record was not found.
var ErrNoRecord = errors.New("models: no matching record found")

// Snippet definition of a snippet record in the database.
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}
