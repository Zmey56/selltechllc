package pkg

import "sync"

var (
	UpdatingFlag bool
	DBMutex      sync.Mutex
)
