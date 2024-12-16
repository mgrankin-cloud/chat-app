package utils

import (
	"sync"
)

var (
	serverList     = []string{}
	LBAddress      = ":3000"
	selectedServer = 0
	mu             sync.Mutex
)
