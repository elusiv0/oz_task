package util

import (
	"errors"
)

func UnwrapError(err error) error {
	next := errors.Unwrap(err)

	for next != nil {
		err = next
		next = errors.Unwrap(err)
	}

	return err
}

type Prid int

func (id *Prid) GenerateId() int {
	genId := *id
	*id++
	return int(genId)
}

func NewPrid() *Prid {
	pr := Prid(1)

	return &pr
}
