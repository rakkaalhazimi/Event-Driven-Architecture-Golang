package database

import (
	"encoding/json"
	"fmt"
	"os"

	"event-driven/models"
	utils "event-driven/utils"
)

type JsonDatabase struct {
	Content []byte
}

func (d *JsonDatabase) GetDocument(key string) (models.Document, error) {
	var index map[string]models.Document
	err := json.Unmarshal(d.Content, &index)
	utils.FailsOnError(err, "Cannot decode json content")

	document := index[key]

	return document, err
}

func (d *JsonDatabase) UpdateDocument(key string, update interface{}) error {
	return nil
}

func NewJsonDatabase(path string) *JsonDatabase {
	content, err := os.ReadFile(path)
	utils.FailsOnError(err, fmt.Sprintf(`Cannot open json file from path '%v'`, path))

	return &JsonDatabase{
		Content: content,
	}
}
