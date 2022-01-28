package usecase

import (
	"context"
	"fmt"
	"time"
	"todoNote/internal/model"
	"todoNote/internal/repo"
	in_memory "todoNote/internal/repo/in-memory"
)

type INoteUsecase interface {
	CreateNote(ctx context.Context, n *model.Note) (model.Id, error)
	FindNote(ctx context.Context, noteId model.Id, userId model.Id, zone model.TimeZone) (*model.Note, error)
	FindAll(ctx context.Context, p FindParams) ([]model.Note, error)
	UpdateNote(ctx context.Context, n *model.Note) error
	RemoveNote(ctx context.Context, noteId, userId model.Id) error
}
var _ INoteUsecase = &NoteUsecase{}

type NoteUsecase struct {
	noteRepo repo.IRepoNote
}

func NewNoteUsecase(r repo.IRepoNote) *NoteUsecase {
	return &NoteUsecase{
		noteRepo: r,
	}
}

//TODO user id from context
func(u *NoteUsecase) CreateNote(ctx context.Context, n *model.Note) (model.Id, error) {
	note := u.prepareNoteDate(n)

	id, err := u.noteRepo.Insert(ctx, note)
	if err != nil {
		return 0, fmt.Errorf("create note: %w", err)
	}
	n.Id = id

	return id, nil
}

func(u *NoteUsecase) FindNote(ctx context.Context, noteId model.Id, userId model.Id, zone model.TimeZone) (*model.Note, error) {
	n, err := u.noteRepo.GetById(ctx, noteId)
	if _, ok := err.(in_memory.NoSuchElementError); ok {
		return nil, NewNoteNotFoundError(noteId, userId)
	}
	if err != nil {
		return nil, fmt.Errorf("find note %w", err)
	}

	if n.UserId != userId {
		return nil, NewNoteNotFoundError(noteId, userId)
	}

	n.Date = Convert(n.Date, zone)

	return &n, err
}

type FindParams struct {
	Filter repo.NoteFilter
	Zone model.TimeZone
}

func(u *NoteUsecase) FindAll(ctx context.Context, p FindParams) ([]model.Note, error) {
	if p.Filter.TakeFrom != nil {
		t := Convert(*p.Filter.TakeFrom, model.UTC)
		p.Filter.TakeFrom = &t
	}

	notes, err := u.noteRepo.GetAllOffset(ctx, p.Filter)
	if err != nil {
		return notes, fmt.Errorf("find all: %w", err)
	}

	return u.mapZone(notes, p.Zone), nil
}

func(u *NoteUsecase) UpdateNote(ctx context.Context, n *model.Note) error {
	note, err := u.noteRepo.GetById(ctx, n.Id)
	if _, ok := err.(in_memory.NoSuchElementError); ok {
		return NewNoteNotFoundError(n.Id, n.UserId)
	}

	if err != nil {
		return fmt.Errorf("update note: %w", err)
	}

	if n.UserId != note.UserId {
		return NewNoteNotFoundError(n.Id, n.UserId)
	}

	updated := u.provideNoteUpdate(&note, n)
	updated = u.prepareNoteDate(updated)

	if err := u.noteRepo.Update(ctx, updated); err != nil {
		return fmt.Errorf("update note: %w", err)
	}

	return nil
}

func(u *NoteUsecase) RemoveNote(ctx context.Context, noteId, userId model.Id) error {
	n, err := u.noteRepo.GetById(ctx, noteId)
	if err != nil {
		if _, ok := err.(in_memory.NoSuchElementError); ok {
			return NewNoteNotFoundError(noteId, userId)
		}

		return fmt.Errorf("remove note: %w", err)
	}

	if n.UserId != userId {
		return NewNoteNotFoundError(noteId, userId)
	}

	if err := u.noteRepo.Delete(ctx, noteId); err != nil {
		return fmt.Errorf("remove note: %w", err)
	}

	return nil
}

func(u *NoteUsecase) prepareNoteDate(n *model.Note) *model.Note {
	if n.Date.IsZero() {
		n.Date = time.Now()
	}

	n.Date = Convert(n.Date, model.UTC)
	return n
}

func(u *NoteUsecase) provideNoteUpdate(old, new *model.Note) *model.Note {
	if new.Text != old.Text && new.Text != ""{
		old.Text = new.Text
	}

	if new.Title != old.Title && new.Title != "" {
		old.Title = new.Title
	}

	if new.IsFinished != old.IsFinished && !old.IsFinished {
		old.IsFinished = new.IsFinished
	}

	if new.Date != old.Date && !new.Date.IsZero(){
		old.Date = new.Date
	}

	return old
}

func(u *NoteUsecase) mapZone(notes []model.Note, zone model.TimeZone) []model.Note {
	for i, n := range notes {
		n.Date = Convert(n.Date, zone)
		notes[i] = n
	}

	return notes
}
