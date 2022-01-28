package usecase

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"todoNote/internal/model"
	"todoNote/internal/repo"
	in_memory "todoNote/internal/repo/in-memory"
	"todoNote/internal/usecase/mocks"
)

//go:generate mockgen -package=mocks -destination=mocks/notes.go todoNote/internal/repo IRepoNote

func TestNoteUsecase_CreateNote(t *testing.T) {
	tts := []struct{
		in model.Note
		want model.Note
		errMsg string
		outErr error
	}{
		{
			in: *model.NewNote(0, 2, "new note", "note", time.Now(), false),
			want: *model.NewNote(0, 2, "new note", "note", time.Now().UTC(), false),
		},
		{
			in: *model.NewNote(0, 2, "my new note", "my note", time.Now(), true),
			want: *model.NewNote(0, 2, "my new note", "my note", time.Now().UTC(), true),
		},
	}

	for _, tt := range tts {
		t.Run("", func(t *testing.T) {
			ctr := gomock.NewController(t)
			defer ctr.Finish()

			mockRepo := mocks.NewMockIRepoNote(ctr)
			mockRepo.EXPECT().Insert(gomock.Any(), gomock.Any()).
				Return(tt.in.Id, tt.outErr).
				Do(func(ctx context.Context, n *model.Note){
					assert.Equal(t, tt.want, *n)
				})

			uc := NewNoteUsecase(mockRepo)

			id, err := uc.CreateNote(context.Background(), &tt.in)
			assert.Equal(t, tt.want.Id, id)
			if err != nil {
				assert.Equal(t, tt.errMsg, err.Error())
			}
		})
	}
}

func TestNoteUsecase_FindNote(t *testing.T) {
	tts := []struct {
		noteId model.Id
		userId model.Id
		zone model.TimeZone
		errMsg string
		testName string
		out model.Note
		outError error
	} {
		{1, 2, model.UTCm4, "", "right id: conversion test",
			model.Note{
			Id: 1,
			UserId: 2,
			Date: time.Now(),
		}, nil},
		{1, 2, model.UTCp4, "", "right id: conversion test",
			model.Note{
			Id: 1,
			UserId: 2,
			Date: time.Now(),
		}, nil},
		{1, 3, model.UTCp4, "no such note found (id: 1) for user.go (id: 3)", "user don't have this note",
			model.Note{
				Id: 1,
				UserId: 2,
				Date: time.Now(),
			}, nil},
		{3, 3, model.UTCp4, NewNoteNotFoundError(3, 3).Error(), "no such note",
			model.Note{},  in_memory.NewNoSuchElementError(3)},
	}


	for _, tt := range tts {
		tn := fmt.Sprintf("note: %v, user: %v, zone: %v,  %v", tt.noteId, tt.userId, tt.zone, tt.testName)
		t.Run(tn, func(t *testing.T){
			mockCtr := gomock.NewController(t)
			defer mockCtr.Finish()

			mockNoteRepo := mocks.NewMockIRepoNote(mockCtr)
			mockNoteRepo.EXPECT().GetById(context.Background(), tt.noteId).
				Return(tt.out, tt.outError)

			uc := NewNoteUsecase(mockNoteRepo)

			got, err := uc.FindNote(context.Background(), tt.noteId, tt.userId, tt.zone)
			if err != nil {
				assert.Equal(t, tt.errMsg, err.Error(), tt.testName)
			}

			if tt.errMsg == "" {
				date := Convert(got.Date, tt.zone)
				assert.Equal(t, date, got.Date, tt.testName)
			}
		})
	}
}

func TestNoteUsecase_FindAll(t *testing.T) {
	now := time.Now().UTC()

	filter := repo.NoteFilter{}
	tts := []struct{
		userId model.Id
		filter FindParams
		out []model.Note
		repoOut []model.Note
		repoErr error
		hasErr bool
	} {
		{
			userId: 1,
			filter: FindParams{Zone: model.UTCp11, Filter: filter},
			out: []model.Note{},
			hasErr: false,
			repoOut: []model.Note{},
			repoErr: nil,
		},
		{
			userId: 2,
			filter: FindParams{Zone: model.UTCp9, Filter: filter},
			out: []model.Note{
				{Id: 1, UserId: 1, Date: Convert(now, model.UTCp9)},
				{Id: 2, UserId: 1, Date: Convert(now, model.UTCp9)},
				{Id: 3, UserId: 1, Date: Convert(now, model.UTCp9)}},
			hasErr: false,
			repoOut: []model.Note{
				{Id: 1, UserId: 1, Date: now},
				{Id: 2, UserId: 1, Date: now},
				{Id: 3, UserId: 1, Date: now},
			},
			repoErr: nil,
		},
		{
			userId: 3,
			filter: FindParams{Zone: model.UTCp5, Filter: filter},
			out: nil,
			hasErr: true,
			repoOut: nil,
			repoErr: fmt.Errorf("some error"),
		},
	}

	for _, tt := range tts {
		tt.filter.Filter.UserId = tt.userId
		tn := fmt.Sprintf("userId: %v, zone: %v", tt.userId, tt.filter.Zone)
		t.Run(tn, func(t *testing.T){
			ctr := gomock.NewController(t)
			defer ctr.Finish()
			mockRepo := mocks.NewMockIRepoNote(ctr)
			mockRepo.EXPECT().GetAllOffset(context.Background(), tt.filter.Filter).
				Return(tt.repoOut, tt.repoErr)

			uc := NewNoteUsecase(mockRepo)

			got, err := uc.FindAll(context.Background(), tt.filter)
			assert.Equal(t, tt.out, got)
			assert.Equal(t, tt.hasErr, err != nil)
			if err != nil {
				assert.Equal(t, fmt.Sprintf("find all: %v", tt.repoErr.Error()),err.Error() )
			}
		})
	}
}

func TestNoteUsecase_UpdateNote(t *testing.T) {
	oldTitle := "old title"
	oldText := "old text"
	newTitle := "new title"
	newText := "new text"
	date := time.Now().UTC()
	newDate := time.Now().Add(2*time.Hour)
	wantDate := newDate.UTC()

	tts := []struct {
		in model.Note
		getByIdNote model.Note
		getByIdErr error
		updateErr error
		want model.Note
		desc string
	} {
		{
			in: model.Note{Id: 1, UserId: 1, Title: newTitle},
			getByIdNote: model.Note{Id: 1, UserId: 1, Title: oldTitle, Text: oldText, IsFinished: true, Date: date},
			getByIdErr: nil,
			updateErr: nil,
			want: model.Note{Id: 1, UserId: 1, Title: newTitle, Text: oldText, IsFinished: true, Date: date},
			desc: "new title",
		},
		{
			in: model.Note{Id: 1, UserId: 1, Text: newText},
			getByIdNote: model.Note{Id: 1, UserId: 1, Title: oldTitle, Text: oldText, IsFinished: true, Date: date},
			getByIdErr: nil,
			updateErr: nil,
			want: model.Note{Id: 1, UserId: 1, Title: oldTitle, Text: newText, IsFinished: true, Date: date},
			desc: "new text",
		},
		{
			in: model.Note{Id: 1, UserId: 1, IsFinished: true},
			getByIdNote: model.Note{Id: 1, UserId: 1, Title: oldTitle, Text: oldText, IsFinished: false, Date: date},
			getByIdErr: nil,
			updateErr: nil,
			want: model.Note{Id: 1, UserId: 1, Title: oldTitle, Text: oldText, IsFinished: true, Date: date},
			desc: "new isFinished",
		},
		{
			in: model.Note{Id: 1, UserId: 1, Date: newDate},
			getByIdNote: model.Note{Id: 1, UserId: 1, Title: oldTitle, Text: oldText, IsFinished: true, Date: date},
			getByIdErr: nil,
			updateErr: nil,
			want: model.Note{Id: 1, UserId: 1, Title: oldTitle, Text: oldText, IsFinished: true, Date: wantDate},
			desc: "new date",
		},
		{
			in: model.Note{Id: 1, UserId: 1, Title: newTitle, Text: newText, IsFinished: true, Date: newDate},
			getByIdNote: model.Note{Id: 1, UserId: 1, Title: oldTitle, Text: oldText, IsFinished: false, Date: date},
			getByIdErr: nil,
			updateErr: nil,
			want: model.Note{Id: 1, UserId: 1,Title: newTitle, Text: newText, IsFinished: true, Date: wantDate},
			desc: "new all",
		},
		{
			in: model.Note{Id: 1, UserId: 1},
			getByIdNote: model.Note{Id: 1, UserId: 1, Title: oldTitle, Text: oldText, IsFinished: true, Date: date},
			getByIdErr: nil,
			updateErr: nil,
			want: model.Note{Id: 1, UserId: 1, Title: oldTitle, Text: oldText, IsFinished: true, Date: date},
			desc: "nothing new",
		},
		{
			in: model.Note{Id: 1, UserId: 2},
			getByIdNote: model.Note{Id: 1, UserId: 1, Title: oldTitle, Text: oldText, IsFinished: true, Date: date},
			getByIdErr: nil,
			updateErr: NewNoteNotFoundError(1, 2),
			want: model.Note{Id: 1, UserId: 2, Title: oldTitle, Text: oldText, IsFinished: true, Date: date},
			desc: "wrong user id",
		},
		{
			in: model.Note{Id: 1, UserId: 1},
			getByIdNote: model.Note{},
			getByIdErr: in_memory.NewNoSuchElementError(1),
			updateErr: NewNoteNotFoundError(1, 1),
			//want: model.Note{Id: 1, UserId: 1, Title: oldTitle, Text: oldText, IsFinished: true, Date: date},
			desc: "wrong note id",
		},
	}

	for _, tt := range tts {
		t.Run(tt.desc, func(t *testing.T) {
			ctr := gomock.NewController(t)
			defer ctr.Finish()
			mockRepo := mocks.NewMockIRepoNote(ctr)

			mockRepo.EXPECT().GetById(gomock.Any(), gomock.Any()).
				Return(tt.getByIdNote, tt.getByIdErr)

			if tt.getByIdErr == nil && tt.updateErr == nil {
				mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).
					Return(tt.updateErr).
					Do(func(ctx context.Context, note *model.Note) {
						assert.Equal(t, tt.want, *note)
					})
			}

			uc := NewNoteUsecase(mockRepo)

			err := uc.UpdateNote(context.Background(), &tt.in)
			if err != nil {
				assert.Equal(t, tt.updateErr.Error(), err.Error())
			}
		})
	}
}

func TestNoteUsecase_RemoveNote(t *testing.T) {
	tts := []struct{
		inNoteId model.Id
		inUserId model.Id
		outErr error
		wantErr error
		wantNoteId model.Id
		storedNote model.Note
	} {
		{
			inNoteId: 1,
			inUserId: 1,
			outErr: nil,
			wantErr: nil,
			wantNoteId: 1,
			storedNote: model.Note{Id: 1, UserId: 1},
		},
		{
			inNoteId: 2,
			inUserId: 1,
			outErr: nil,
			wantErr: nil,
			wantNoteId: 2,
			storedNote: model.Note{Id: 2, UserId: 1},
		},
		{
			inNoteId: 2,
			inUserId: 2,
			outErr: nil,
			wantErr: NewNoteNotFoundError(2,2),
			wantNoteId: 2,
			storedNote: model.Note{Id: 2, UserId: 1},
		},
		{
			inNoteId: 2,
			inUserId: 1,
			outErr: nil,
			wantErr: NewNoteNotFoundError(2,1),
			wantNoteId: 2,
			storedNote: model.Note{Id: 2, UserId: 2},
		},
	}

	for _, tt := range tts {
		t.Run("", func(t *testing.T) {
			ctr := gomock.NewController(t)
			defer ctr.Finish()
			mockRepo := mocks.NewMockIRepoNote(ctr)
			if tt.wantErr == nil {
				mockRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).
					Return(tt.outErr).
					Do(func(ctx context.Context, noteId model.Id) {
						assert.Equal(t, tt.wantNoteId, noteId)
					})
			}
			mockRepo.EXPECT().GetById(gomock.Any(), gomock.Any()).
				Return(tt.storedNote, nil)

			u := NewNoteUsecase(mockRepo)

			err := u.RemoveNote(context.Background(), tt.inNoteId, tt.inUserId)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestNoteUsecase_PrepareNoteDate(t *testing.T) {
	uc := NoteUsecase{}

	tts := []struct{
		in model.Note
	} {
		{model.Note{}},
		{model.Note{Date: time.Now().Local()}},
		{model.Note{Date: Convert(time.Now(), model.UTCp9)}},
	}

	for _, tt := range tts {
		got := uc.prepareNoteDate(&tt.in)
		assert.Equal(t, got.Date.IsZero(), false)
		_, off := got.Date.Zone()
		assert.Equal(t, off, 0)
	}
}

func TestNoteUsecase_ProvideNoteUpdate(t *testing.T) {
	uc := NoteUsecase{}
	dt := time.Now()
	dtNew := dt.Add(time.Minute)
	tts := []struct{
		old model.Note
		new model.Note
		out model.Note
	} {
		{
			old: model.Note{Id: 10, UserId: 11, Title: "old title", Text: "old text", IsFinished: false, Date: dt},
			new: model.Note{Id: 13, UserId: 14, Title: "new title", Text: "new text", IsFinished: true, Date: dtNew},
			out: model.Note{Id: 10, UserId: 11, Title: "new title", Text: "new text", IsFinished: true, Date: dtNew},
		},
		{
			old: model.Note{Id: 10, UserId: 11, Title: "old title", Text: "old text", IsFinished: false, Date: dt},
			new: model.Note{Id: 13, UserId: 14, Title: "old title", Text: "old text", IsFinished: false, Date: dt},
			out: model.Note{Id: 10, UserId: 11, Title: "old title", Text: "old text", IsFinished: false, Date: dt},
		},
	}
	for _, tt := range tts {
		assert.Equal(t, tt.out, *uc.provideNoteUpdate(&tt.old, &tt.new))
	}
}

func TestNoteUsecase_MapZone(t *testing.T) {
	uc := NoteUsecase{}
	date1 := time.Now().UTC()
	date2 := time.Now().Add(5*time.Hour).UTC()
	date3 := date1.Add(22*time.Minute).UTC()

	tts := []struct{
		in []model.Note
		out []model.Note
		zone model.TimeZone
	} {
		{
			in: []model.Note{{Date: date1}, {Date: date2}, {Date: date3}},
			out: []model.Note{{Date: date1}, {Date: date2}, {Date: date3}},
			zone: model.UTC,
		},
		{
			in: []model.Note{{Date: date1}, {Date: date2}, {Date: date3}},
			out: []model.Note{{Date: Convert(date1, model.UTCp9)}, {Date: Convert(date2, model.UTCp9)}, {Date: Convert(date3, model.UTCp9)}},
			zone: model.UTCp9,
		},
	}
	for _, tt := range tts {
		assert.Equal(t, tt.out, uc.mapZone(tt.in, tt.zone))
	}
}

func BenchmarkNoteUsecase_MapZone(b *testing.B) {
	uc := NoteUsecase{}
	notes := make([]model.Note, 0, 100)
	for i := 0; i < 100; i++ {
		notes = append(notes,
			model.Note{
				Id: model.Id(i),
				Date: time.Now().UTC(),
			})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		uc.mapZone(notes, model.UTCp3)
	}
}

