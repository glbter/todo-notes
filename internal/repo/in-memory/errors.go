package in_memory

import (
	"fmt"
	"todoNote/internal/model"
)

type NoSuchElementError struct {
	id int64
}

func NewNoSuchElementError(id model.Id) NoSuchElementError {
	return NoSuchElementError{id: id}
}

func (err NoSuchElementError) Error() string {
	return fmt.Sprintf("no such element with id %v", err.id)
}

type NoSuchNameError struct {
	name string
}

func NewNoSuchNameError(name string ) NoSuchNameError {
	return NoSuchNameError{name: name}
}

func (err NoSuchNameError) Error() string {
	return fmt.Sprintf("no such element with name %v", err.name)
}

