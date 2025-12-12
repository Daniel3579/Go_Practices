package core

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("not found")

func Now() time.Time {
	return time.Now().UTC()
}
