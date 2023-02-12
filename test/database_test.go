package test

import (
	"fmt"
	"testing"

	database "event-driven/database"
)

func TestReadJson(t *testing.T) {
	path := "../data/documents.json"
	db := database.NewJsonDatabase(path)

	document, _ := db.GetDocument("0")

	fmt.Println(document)
}
