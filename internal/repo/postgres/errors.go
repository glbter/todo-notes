package postgres

import "fmt"

const (
	users      = "users:"
	notes = "notes:"
	select_sql = "select:"
	insert = "insert:"
	delete_sql = "delete_sql:"
	update     = "update:"

)

type UserExistsError struct {
	name string
}

func (e UserExistsError) Error() string {
	return fmt.Sprintf("user (%v) already exists", e.name)
}

var rowsAffectedNotOne = fmt.Errorf("rows affected != 1")

func NewUsersError(method string, err error) error {
	return fmt.Errorf("%v %v %w", users, method, err)
}

func NewNotesError(method string, err error) error {
	return fmt.Errorf("%v %v %w", notes, method, err)
}
