package main

import (
	"sync"

	"github.com/xtaci/kcptun/generic"
)

var (
	multiPort           *generic.MultiPort
	multiPortParseError error
	multiPortOnce       sync.Once
)
