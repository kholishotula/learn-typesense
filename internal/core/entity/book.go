package entity

import (
	"encoding/json"
	"fmt"
)

type Book struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Year    int    `json:"year"`
	Summary string `json:"summary"`
}

type CreateInput struct {
	Title   string `json:"title"`
	Author  string `json:"author"`
	Year    int    `json:"year"`
	Summary string `json:"summary"`
}

func (b *Book) EncodeJSON() string {
	bytes, _ := json.Marshal(b)
	return string(bytes)
}

func (b *Book) DecodeJSON(src string) error {
	err := json.Unmarshal([]byte(src), b)
	if err != nil {
		return fmt.Errorf("unable to decode source due: %w", err)
	}
	return nil
}
