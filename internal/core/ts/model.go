package ts

import (
	"encoding/json"
	"fmt"
)

type Doc interface {
	EncodeJSON() string
	DecodeJSON(src string) error
}

func (d *MapDoc) EncodeJSON() string {
	bytes, _ := json.Marshal(d.Source)
	return string(bytes)
}

func (d *MapDoc) DecodeJSON(src string) error {
	err := json.Unmarshal([]byte(src), &d.Source)
	if err != nil {
		return fmt.Errorf("unable to decode source due: %w", err)
	}
	return nil
}

type MapDoc struct {
	Source map[string]interface{}
}

type IndexInput struct {
	CollectionName string
	Doc            Doc
}

type RetrieveInput struct {
	CollectionName string
	ID             string
}

type RetrieveOutput struct {
	CollectionName string
	Doc            Doc
}

type SearchInput struct {
	CollectionName string
	Query          string
	QueryBy        string
	Limit          int
	Page           int
}

type SearchOutput struct {
	Found int   `json:"found"`
	Hits  []Doc `json:"hits"`
}
