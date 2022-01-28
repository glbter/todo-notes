package usecase

import "fmt"

type ElemNotFound struct {
	ElemId int64
	UserId int64
	TypeName string
}

func NewElemNotFoundError(name string, eId, uId int64) *ElemNotFound {
	return &ElemNotFound{ElemId: eId, UserId: uId, TypeName: name}
}

const noteType = "note"
func NewNoteNotFoundError(eId, uId int64) *ElemNotFound {
	return NewElemNotFoundError(noteType, eId, uId)
}

func(e ElemNotFound) Error() string {
	return fmt.Sprintf("no such %v found (id: %v) for user.go (id: %v)", e.TypeName, e.ElemId, e.UserId)
}

