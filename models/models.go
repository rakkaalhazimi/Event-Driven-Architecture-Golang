package models

type QueueMessages struct {
	Key     string `json:"key" bson:"key"`
	Service string `json:"service" bson:"service"`
}

type Document struct {
	Text           string `json:"text" bson:"text"`
	CleanText      string `json:"clean_text" bson:"clean_text"`
	SynthesizeType string `json:"synthesize_type" bson:"synthesize_type"`
}
