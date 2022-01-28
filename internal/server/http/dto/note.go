package dto

import "time"

type NewNote struct {
	Title string `json:"title"`
	Text string `json:"text"`
	Date time.Time `json:"date"`
}

type Note struct {
	Id string `json:"id"`
	Title string `json:"title,omitempty"`
	Text string `json:"text,omitempty"`
	Date time.Time `json:"date,omitempty"`
	IsFinished bool `json:"is_finished,omitempty"`
}

type Notes = []Note

type NoteUpdate struct {
	Title string `json:"title,omitempty"`
	Text string `json:"text,omitempty"`
	Date time.Time `json:"date,omitempty"`
	IsFinished bool `json:"is_finished,omitempty"`
}

