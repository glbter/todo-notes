

package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"todoNote/internal/model"
	"todoNote/internal/repo"
	in_memory "todoNote/internal/repo/in-memory"
)


//func notesInit(t *testing.T) {
//	usersInit(t)
//
//	c := connect(t)
//	c.Exec(context.Background(), `insert into notes
//	(id, user_id, title, text, date, is_finished)
//	values
//	(1, 1, 'title', 'text', '2021-10-13 10:20:08.392115', false),
//	(2, 1, 'title', 'text', '2021-10-12 10:20:08.392115', false),
//	(3, 1, 'title', 'text', '2021-10-13 10:20:08.392115', false),
//	(4, 1, 'title', 'text', '2021-10-13 10:20:08.392115', true),
//	(5, 1, 'title', 'text', '2021-10-12 10:20:08.392115', true),
//	(6, 1, 'title', 'text', '2021-10-13 10:20:08.392115', true),
//
//	(7, 2, 'title', 'text', '2021-10-12 10:20:08.392115', true),
//	(8, 2, 'title', 'text', '2021-10-13 10:20:08.392115', true),
//	(9, 2, 'title', 'text', '2021-10-12 10:20:08.392115', false);`)
//	c.Close(context.Background())
//}

func areNotesEqual(t *testing.T, n1, n2 model.Note) {
	assert.Equal(t, n1.Id, n2.Id)
	assert.Equal(t, n1.Text, n2.Text)
	assert.Equal(t, n1.Title, n2.Title)
	assert.Equal(t, n1.UserId, n2.UserId)
	assert.Equal(t, n1.IsFinished, n2.IsFinished)
	assert.Equal(t, n1.Date.Format(time.RFC3339), n2.Date.Format(time.RFC3339))
}

func getNote(t *testing.T, conn *pgx.Conn, id model.Id) model.Note {
	var out model.Note
	err := conn.QueryRow(context.Background(),
		`select id, title, text, user_id, date, is_finished from notes where id = $1`,
		id).Scan(
		&out.Id,
		&out.Title,
		&out.Text,
		&out.UserId,
		&out.Date,
		&out.IsFinished)
	assert.Nil(t, err)

	return out
}

func connect(t *testing.T) *pgx.Conn {
	ctx := context.Background()
	c, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Ping(ctx); err != nil {
		t.Fatal(err)
	}

	return c
}



func TestRepoNote_Insert(t *testing.T) {
	//t.Run("init", func(t *testing.T) {
	//	usersInit(t)
	//})

	tts := []struct{
		caseName string
		hasError bool
	} {
		{
			caseName: "success",
			hasError: false,
		},
	}
	for _, tt := range tts {
		t.Run(tt.caseName, func(t *testing.T) {
			conn := connect(t)
			defer conn.Close(context.Background())

			r := NewRepoNote(conn)

			in := model.NewNote(0, model.Id(1), "new_title", "new_text", time.Now().UTC(), false)
			id, err := r.Insert(context.Background(), in)
			in.Id = id
			assert.Equal(t, tt.hasError, err != nil)
			assert.NotEqual(t, 0, id)

			if !tt.hasError {
				out := getNote(t, conn, id)
				areNotesEqual(t, *in, out)
			}
		})
	}

	//t.Run("clean up", func(t *testing.T) {
	//	usersClean(t)
	//})
}

func TestRepoNote_GetById(t *testing.T) {
	//t.Run("init", func(t *testing.T) {
	//	notesInit(t)
	//})

	tts := []struct{
		caseName string
		hasError bool
	} {
		{
			caseName: "success",
			hasError: false,
		},
		{
			caseName: "fail",
			hasError: true,
		},
	}
	for _, tt := range tts {
		t.Run(tt.caseName, func(t *testing.T) {
			conn := connect(t)
			defer conn.Close(context.Background())

			r := NewRepoNote(conn)

			//for id := 1; id < 10; id++ {
			id := 7
			var in model.Note
			if !tt.hasError {
				in = getNote(t, conn, model.Id(id))
			}

			out, err := r.GetById(context.Background(), model.Id(id))
			assert.Equal(t, tt.hasError, err != nil)
			assert.NotEqual(t, 0, id)

			if !tt.hasError {
				areNotesEqual(t, out, in)
			}
			//}

			//usersClean(t)
			conn.Exec(context.Background(), `delete from notes where id = $1`, id)
		})
	}
}

func TestRepoNote_Update(t *testing.T) {
	//t.Run("init", func(t *testing.T) {
	//	notesInit(t)
	//})

	tts := []struct{
		caseName string
		hasError bool
	} {
		{
			caseName: "success",
			hasError: false,
		},
	}
	for _, tt := range tts {
		t.Run(tt.caseName, func(t *testing.T) {
			conn := connect(t)
			defer conn.Close(context.Background())

			r := NewRepoNote(conn)

			//for id := 1; id < 10; id++ {
			id := 8
			var in model.Note
			if !tt.hasError {
				in = getNote(t, conn, model.Id(id))
			}

			in.Text = "txt"
			in.Title = "ttl"
			in.IsFinished = true
			in.Date = time.Now().UTC()
			tmpUid := in.UserId
			in.UserId = 0

			err := r.Update(context.Background(), &in)
			assert.Equal(t, tt.hasError, err != nil)
			assert.NotEqual(t, 0, id)

			if !tt.hasError {
				in.UserId = tmpUid
				out := getNote(t, conn, model.Id(id))
				areNotesEqual(t, out, in)
			}
			//}

			conn.Exec(context.Background(), `delete from notes where id = $1`, id)
		})
	}
}

func TestRepoNote_Delete(t *testing.T) {
	//t.Run("init", func(t *testing.T) {
	//	notesInit(t)
	//})

	tts := []struct{
		caseName string
		hasError bool
	} {
		{
			caseName: "success",
			hasError: false,
		},
		{
			caseName: "fail",
			hasError: true,
		},
	}
	for _, tt := range tts {
		t.Run(tt.caseName, func(t *testing.T) {
			conn := connect(t)
			defer conn.Close(context.Background())

			r := NewRepoNote(conn)

			//for id := 1; id < 10; id++ {
				id := 9
				err := r.Delete(context.Background(), model.Id(id))
				assert.Equal(t, tt.hasError, err != nil)
			//}

			//usersClean(t)
		})
	}
}

func TestRepoNote_Insert_GetById_Delete(t *testing.T) {
	t.Run("", func(t *testing.T) {
		conn := connect(t)
		defer conn.Close(context.Background())
		rn := NewRepoNote(conn)
		ru := NewRepoUser(conn)

		u := model.NewUser(0, "user", []byte("user"), model.UTC)

		uId, err := ru.Insert(context.Background(), u)
		if err != nil {
			t.Fatal(err)
		}
		defer ru.Delete(context.Background(), uId)

		in := model.NewNote(0, uId, "title", "text", time.Now().UTC(), false)
		id, err := rn.Insert(context.Background(), in)
		assert.Nil(t, err)
		assert.NotEqual(t, 0, id)

		out, err := rn.GetById(context.Background(), id)
		assert.Nil(t, err)
		assert.Equal(t, in.Text, out.Text)
		assert.Equal(t, in.Title, out.Title)
		assert.Equal(t, in.UserId, out.UserId)
		assert.Equal(t, in.IsFinished, out.IsFinished)
		assert.Equal(t, in.Date.Format(time.RFC3339), out.Date.Format(time.RFC3339))

		err = rn.Delete(context.Background(), id)
		assert.Nil(t, err)

		_, err = rn.GetById(context.Background(), id)
		assert.NotNil(t, err)

		assert.Equal(t, true, errors.As(err, &in_memory.NoSuchElementError{}))
	})
}

func TestRepoNote_Update2(t *testing.T) {
	tts := []struct{
		up int64
		caseName string
	} {
		{0, "success: userIds are same"},
		{1, "success: userIds are not same"},
	}
	for _, tt := range tts {
		t.Run("success", func(t *testing.T) {
			conn := connect(t)
			defer conn.Close(context.Background())
			rn := NewRepoNote(conn)
			ru := NewRepoUser(conn)

			u := model.NewUser(0, "user", []byte("user"), model.UTC)

			uId, err := ru.Insert(context.Background(), u)
			if err != nil {
				t.Fatal(err)
			}
			defer ru.Delete(context.Background(), uId)

			in := model.NewNote(0, uId, "title", "text", time.Now().UTC(), false)
			id, _ := rn.Insert(context.Background(), in)

			upd := model.NewNote(id, uId+tt.up, "new title", "new text", time.Now().UTC().Add(2 * time.Second), true)
			err = rn.Update(context.Background(), upd)
			assert.Nil(t, err)

			out, err := rn.GetById(context.Background(), id)
			assert.Nil(t, err)
			assert.Equal(t, upd.Text, out.Text)
			assert.Equal(t, upd.Title, out.Title)
			assert.Equal(t, upd.UserId-tt.up, out.UserId)
			assert.Equal(t, upd.IsFinished, out.IsFinished)
			assert.Equal(t, upd.Date.Format(time.RFC3339), out.Date.Format(time.RFC3339))

			err = rn.Delete(context.Background(), id)
			assert.Nil(t, err)
		})
	}
}


func TestRepoNote_GetAllOffset(t *testing.T) {
	t.Run("two users seven notes", func(t *testing.T) {
		conn := connect(t)
		defer conn.Close(context.Background())
		rn := NewRepoNote(conn)
		ru := NewRepoUser(conn)

		u := model.NewUser(0, "user", []byte("user"), model.UTC)
		u2 := model.NewUser(0, "user2", []byte("user"), model.UTC)

		uId, _ := ru.Insert(context.Background(), u)
		defer ru.Delete(context.Background(), uId)
		uId2, _ := ru.Insert(context.Background(), u2)
		defer ru.Delete(context.Background(), uId2)

		now := time.Now().UTC()
		later := now.Add(2 * time.Hour)

		notes := []*model.Note{
			model.NewNote(0, uId, "title1", "text1", now, false),
			model.NewNote(0, uId, "title2", "text2", now, false),
			model.NewNote(0, uId, "title3", "text3", now, false),
			model.NewNote(0, uId, "title4", "text4", later, true),
			model.NewNote(0, uId, "title5", "text5", later, true),
			model.NewNote(0, uId2, "title6", "text6", later, true),
			model.NewNote(0, uId2, "title7", "text7", later, true),
		}

		ids := make([]model.Id, 0, len(notes))
		for _, note := range notes {
			id, _ := rn.Insert(context.Background(), note)
			ids = append(ids, id)

			defer func(id model.Id) {
				rn.Delete(context.Background(), id)
			}(id)
		}

		t.Run(fmt.Sprintf("by userId 1"), func(t *testing.T) {
			filter := repo.NoteFilter{UserId: uId}
			out, err := rn.GetAllOffset(context.Background(), filter)
			assert.Nil(t, err)
			assert.Equal(t, 5, len(out))
		})

		t.Run(fmt.Sprintf("by userId 2"), func(t *testing.T) {
			filter := repo.NoteFilter{UserId: uId2}
			out, err := rn.GetAllOffset(context.Background(), filter)
			assert.Nil(t, err)
			assert.Equal(t, 2, len(out))
		})

		t.Run(fmt.Sprintf("by userId 2, not finished"), func(t *testing.T) {
			notFinish := false
			filter := repo.NoteFilter{UserId: uId2, IsFinished: &notFinish}
			out, err := rn.GetAllOffset(context.Background(), filter)
			assert.Nil(t, err)
			assert.Equal(t, 0, len(out))
		})

		t.Run(fmt.Sprintf("by userId 1, not finished"), func(t *testing.T) {
			notFinish := false
			filter := repo.NoteFilter{UserId: uId, IsFinished: &notFinish}
			out, err := rn.GetAllOffset(context.Background(), filter)
			assert.Nil(t, err)
			assert.Equal(t, len(out), 3)
		})

		t.Run(fmt.Sprintf("by userId 1, later"), func(t *testing.T) {
			from := now.Add(time.Hour)
			filter := repo.NoteFilter{UserId: uId, TakeFrom: &from}
			out, err := rn.GetAllOffset(context.Background(), filter)
			assert.Nil(t, err)
			assert.Equal(t, 2, len(out))
		})

		t.Run(fmt.Sprintf("by userId 1, all time"), func(t *testing.T) {
			from := now.Add(-time.Hour)
			filter := repo.NoteFilter{UserId: uId, TakeFrom: &from}
			out, err := rn.GetAllOffset(context.Background(), filter)
			assert.Nil(t, err)
			assert.Equal(t, len(out), 5)
		})

		t.Run(fmt.Sprintf("by userId 1, offset 1"), func(t *testing.T) {
			off := uint64(ids[0])
			p := repo.PageFilter{Offset: &off}
			filter := repo.NoteFilter{UserId: uId, Page: p}
			out, err := rn.GetAllOffset(context.Background(), filter)
			assert.Nil(t, err)
			assert.Equal(t, len(out), 4)
		})

		t.Run(fmt.Sprintf("by userId 1, offset 2"), func(t *testing.T) {
			off := uint64(ids[1])
			p := repo.PageFilter{Offset:  &off}
			filter := repo.NoteFilter{UserId: uId, Page: p}
			out, err := rn.GetAllOffset(context.Background(), filter)
			assert.Nil(t, err)
			assert.Equal(t, len(out), 3)
		})

		t.Run(fmt.Sprintf("by userId 1, offset 1, limit 3"), func(t *testing.T) {
			off := uint64(ids[0])
			p := repo.PageFilter{
				Offset: &off,
				Limit: repo.GetUIntParamPointer("3"),
			}

			filter := repo.NoteFilter{UserId: uId, Page: p}
			out, err := rn.GetAllOffset(context.Background(), filter)
			assert.Nil(t, err)
			assert.Equal(t, len(out), 3)
		})

		t.Run(fmt.Sprintf("by userId 1, offset 3, limit 3"), func(t *testing.T) {
			off := uint64(ids[2])
			p := repo.PageFilter{
				Offset: &off,
				Limit: repo.GetUIntParamPointer("3"),
			}

			filter := repo.NoteFilter{UserId: uId, Page: p}
			out, err := rn.GetAllOffset(context.Background(), filter)
			assert.Nil(t, err)
			assert.Equal(t, len(out), 2)
		})
	})
}
