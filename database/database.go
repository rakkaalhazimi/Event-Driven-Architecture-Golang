package database

import models "event-driven/models"

type DatabaseObject interface {
	GetDocument(key string) (models.Document, error)
	UpdateDocument(key string, update interface{}) error
}
