package model

import "time"

type Id = int64

type Note struct {
	Id     Id
	UserId Id
	Title  string
	Text string
	Date time.Time
	IsFinished bool
}

func NewNote(id Id, usedId Id, title, text string, date time.Time, isFinished bool) *Note {
	return &Note{
		Id: id,
		UserId: usedId,
		Title: title,
		Text: text,
		Date: date,
		IsFinished: isFinished,
	}
}
