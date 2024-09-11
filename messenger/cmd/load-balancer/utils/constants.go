package utils

import (
	"sync"
)

var (
	serverList     = []string{server.ssoServerAddr, chatSeverAddr, userServerAddr, ntfServerAddr, mediaServerAddr, modelsServerAddr, msgServerAddr}
	LBAddress      = ":3000"
	selectedServer = 0
	mu             sync.Mutex
)
