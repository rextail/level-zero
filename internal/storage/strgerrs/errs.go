package strgerrs

import "errors"

var (
	ErrAlreadyExists    = errors.New("already exists")
	ErrZeroRecordsFound = errors.New("found zero records")
)
