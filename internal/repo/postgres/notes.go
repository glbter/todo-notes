package postgres

import (
	"context"
	"github.com/jackc/pgx/v4"
	"strings"
	"time"
	"todoNote/internal/model"
	"todoNote/internal/repo"
	in_memory "todoNote/internal/repo/in-memory"
)

var _ repo.IRepoNote = RepoNote{}
type RepoNote struct {
	conn *pgx.Conn
}

func NewRepoNote(c *pgx.Conn) *RepoNote {
	return &RepoNote{conn: c}
}

func (r RepoNote) Insert(ctx context.Context, n *model.Note) (model.Id, error) {
	query := `
INSERT INTO notes (user_id, title, text, date, is_finished)
VALUES ($1, $2, $3, $4, $5) RETURNING id;`

	var id model.Id
	err := r.conn.QueryRow(ctx,
		query,
		n.UserId,
		n.Title,
		n.Text,
		n.Date,
		n.IsFinished).
		Scan(&id)

	if err != nil {
		return 0, NewNotesError(insert, err)
	}

	return id, nil
}

func (r RepoNote) GetById(ctx context.Context, noteId model.Id) (model.Note, error) {
	query := `SELECT id, user_id, title, text, date, is_finished FROM notes WHERE id = $1;`
	var note model.Note
	err := r.conn.QueryRow(ctx, query, noteId).Scan(
		&note.Id,
		&note.UserId,
		&note.Title,
		&note.Text,
		&note.Date,
		&note.IsFinished)

	if err != nil {
		isEmpty := strings.Contains(err.Error(), "no rows in result set")
		if isEmpty {
			return model.Note{}, in_memory.NewNoSuchElementError(noteId)
		}

		return model.Note{}, NewNotesError(select_sql, err)
	}

	return note, err
}

func (r RepoNote) GetAllOffset(ctx context.Context, filter repo.NoteFilter) ([]model.Note, error) {
	q := `SELECT id, user_id, title, text, date, is_finished FROM notes
WHERE user_id = $1
AND id > $2 --0
AND (is_finished = $3 OR is_finished = $4) --true and false
AND date > $5 --1970
ORDER BY date 
LIMIT $6;`
	//query := `SELECT id, user_id, title, text, date, is_finished FROM notes WHERE user_id = $1`
	//if filter.IsFinished != nil {
	//	query += fmt.Sprintf(` AND is_finished = %v`, *filter.IsFinished)
	//}
	//if filter.TakeFrom != nil {
	//	dt := filter.TakeFrom.Format(time.RFC3339)
	//	dtT := strings.ReplaceAll(dt, "T", " ")
	//	dtTZ := strings.ReplaceAll(dtT, "Z", "")
 	//	query += fmt.Sprintf(` AND date > '%v'`, dtTZ)
	//}
	//if filter.Page.Offset != nil {
	//	query += fmt.Sprintf(` AND id > %v`, *filter.Page.Offset)
	//}
	//
	//query += ` ORDER BY date`
	//if filter.Page.Limit != nil {
	//	query += fmt.Sprintf(` LIMIT %v`, *filter.Page.Limit)
	//}
	//query += `;`
	//
	//rows, err := r.conn.Query(ctx,
	//	query,
	//	filter.UserId)

	var offset uint64
	is_finished1 := true
	is_finished2 := false
	date := time.Date(1970, 1, 1, 0, 0 ,0 ,0, time.UTC)
	limit := uint64(1000)

	if filter.IsFinished != nil {
		is_finished1 = *filter.IsFinished
		is_finished2 = *filter.IsFinished
	}
	if filter.TakeFrom != nil {
		date = *filter.TakeFrom
	}
	if filter.Page.Offset != nil {
		offset = *filter.Page.Offset
	}

	if filter.Page.Limit != nil {
		limit = *filter.Page.Limit
	}

	rows, err := r.conn.Query(ctx,
		q,
		filter.UserId,
		offset,
		is_finished1,
		is_finished2,
		date,
		limit)

	if err != nil {
		return nil, NewNotesError(select_sql, err)
	}

	res := make([]model.Note, 0)
	defer rows.Close()
	for rows.Next() {
		var row model.Note
		err := rows.Scan(
			&row.Id,
			&row.UserId,
			&row.Title,
			&row.Text,
			&row.Date,
			&row.IsFinished)

		if err != nil {
			return nil, NewNotesError(select_sql, err)
		}

		res = append(res, row)
	}

	if rows.Err() != nil {
		return nil, NewNotesError(select_sql, err)
	}

	return res, err
}

func (r RepoNote) Update(ctx context.Context, n *model.Note) error {
	query := `UPDATE notes SET title = $1, text = $2, date = $3, is_finished = $4 WHERE id = $5;`
	res, err := r.conn.Exec(ctx,
		query,
		n.Title,
		n.Text,
		n.Date,
		n.IsFinished,
		n.Id)

	if err != nil {
		return NewNotesError(update, err)
	}

	if res.RowsAffected() != 1 {
		return NewNotesError(update, rowsAffectedNotOne)
	}

	return nil
}

func (r RepoNote) Delete(ctx context.Context, noteId model.Id) error {
	query := `DELETE FROM notes WHERE id = $1;`
	res, err := r.conn.Exec(ctx,
		query,
		noteId)

	if err != nil {
		return NewNotesError(delete_sql, err)
	}

	if res.RowsAffected() != 1 {
		return NewNotesError(delete_sql, rowsAffectedNotOne)
	}

	return nil
}

