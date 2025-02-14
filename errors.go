package cobrautil

import (
	"errors"
)

var (
	ErrNoHomeDirVar = errors.New("neither $XDG_CONFIG_HOME nor $HOME are defined")
)
