package postgres

import (
	"context"
	"github.com/jackc/pgx/v4"
	"strings"
	"todoNote/internal/model"
	"todoNote/internal/repo"
	in_memory "todoNote/internal/repo/in-memory"
)
var _ repo.IRepoUser = &RepoUser{}

type RepoUser struct {
	conn *pgx.Conn
}

func NewRepoUser(conn *pgx.Conn) repo.IRepoUser {
	return &RepoUser{
		conn: conn,
	}
}

func (r *RepoUser) Insert(ctx context.Context, u *model.User) (model.Id, error) {
	query := `INSERT INTO users (name, password_hash, time_zone) VALUES ($1, $2, $3) RETURNING id;`
	var id model.Id
	err := r.conn.QueryRow(ctx,
		query,
		u.Name,
		u.PasswordHash,
		u.TimeZone).
		Scan(&id)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return id, UserExistsError{name: u.Name}
		}

		return 0, NewUsersError(insert, err)
	}

	return id, nil
}

func (r *RepoUser) GetByUserName(ctx context.Context, name string) (*model.User, error) {
	query := `SELECT id, name, time_zone, password_hash FROM users WHERE name = $1;`
	u, err := r.get(ctx, query, name)

	if err != nil {
		isEmpty := strings.Contains(err.Error(), "no rows in result set")
		if isEmpty {
			return nil, in_memory.NewNoSuchNameError(name)
		}

		return nil, NewUsersError(select_sql, err)
	}

	return u, nil
}

func (r *RepoUser) GetById(ctx context.Context, uId model.Id) (*model.User, error) {
	query := `SELECT id, name, time_zone, password_hash FROM users WHERE id = $1;`
	u, err := r.get(ctx, query, uId)

	if err != nil {
		isEmpty := strings.Contains(err.Error(), "no rows in result set")
		if isEmpty {
			return nil, in_memory.NewNoSuchElementError(uId)
		}

		return nil, NewUsersError(select_sql, err)
	}

	return u, nil
}

func (r *RepoUser) Update(ctx context.Context, u *model.User) error {
	query := `UPDATE users SET name = $1, time_zone = $2 WHERE id = $3;`
	res, err := r.conn.Exec(ctx,
		query,
		u.Name,
		u.TimeZone,
		u.Id)

	if err != nil {
		return NewUsersError(update, err)
	}

	if res.RowsAffected() != 1 {
		return NewUsersError(update, rowsAffectedNotOne)
	}

	return nil
}

func (r *RepoUser) Delete(ctx context.Context, uId model.Id) error {
	query := `DELETE FROM users WHERE id = $1;`
	res, err := r.conn.Exec(ctx,
		query,
		uId)

	if err != nil {
		return NewUsersError(delete_sql, err)
	}

	if res.RowsAffected() != 1 {
		return NewUsersError(delete_sql, rowsAffectedNotOne)
	}

	return nil
}

func (r *RepoUser) get(ctx context.Context, query string, param interface{}) (*model.User, error){
	var usr model.User
	err := r.conn.QueryRow(ctx,
		query,
		param).
		Scan(&usr.Id,
			&usr.Name,
			&usr.TimeZone,
			&usr.PasswordHash)

	return &usr, err
}
